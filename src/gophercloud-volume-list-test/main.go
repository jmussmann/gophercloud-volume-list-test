package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newHttpClient(skipVerifySSL bool) http.Client {
	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerifySSL},
		},
	}
}

var (
	debug         bool
	name          string
	vm            string
	disk          string
	skipVerifySSL bool
)

var rootCmd = &cobra.Command{

	Use: "volume-test",
	RunE: func(cmd *cobra.Command, args []string) error {

		log.SetLevel(log.DebugLevel)

		ctx := context.TODO()

		authOptions, err := openstack.AuthOptionsFromEnv()

		if err != nil {
			panic(err)
		}

		providerClient, err := openstack.NewClient(authOptions.IdentityEndpoint)
		providerClient.HTTPClient = newHttpClient(skipVerifySSL)
		if err != nil {
			panic(err)
		}

		err = openstack.Authenticate(ctx, providerClient, authOptions)

		if err != nil {
			panic(err)
		}

		volumeClient, err := openstack.NewBlockStorageV3(providerClient, gophercloud.EndpointOpts{
			Region: os.Getenv("OS_REGION_NAME"),
		})

		if err != nil {
			panic(err)
		}

		pages, err := volumes.List(volumeClient, volumes.ListOpts{
			Name: name,
		}).AllPages(ctx)

		if err != nil {
			panic(err)
		}

		volumeList, err := volumes.ExtractVolumes(pages)

		if err != nil {
			panic(err)
		}

		if len(volumeList) == 0 {
			log.Debug("Error, volume list empty")
			return nil
		}

		log.Debugf("Info: VolumeList[0].ID: %s", volumeList[0].ID)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&name, "name", "", "volume name")
	rootCmd.MarkPersistentFlagRequired("name")
	rootCmd.PersistentFlags().StringVar(&vm, "vm", "", "VM id in vCenter vm-xxxxx")
	rootCmd.MarkPersistentFlagRequired("vm")
	rootCmd.PersistentFlags().StringVar(&disk, "disk", "", "disk-id xxxx-yyyy")
	rootCmd.MarkPersistentFlagRequired("disk")
	rootCmd.PersistentFlags().BoolVar(&skipVerifySSL, "skip-verify-ssl", false, "disable ssl verification")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
