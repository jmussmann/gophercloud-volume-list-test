package main

import (
	"context"
    "os"
    "fmt"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/gophercloud/gophercloud/v2"
    "github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)

var (
    debug            bool
	name             string
	vm               string
    disk             string
	testing          bool
)

var rootCmd = &cobra.Command{

	Use: "volume-test",
	RunE: func(cmd *cobra.Command, args []string) error {

		log.SetLevel(log.DebugLevel)

		ctx := context.TODO()

	    opts, err := openstack.AuthOptionsFromEnv()
    

        providerClient, err := openstack.AuthenticatedClient(ctx, opts)
        if err != nil {
            panic(err)
        }

        volumeClient, err := openstack.NewBlockStorageV3(providerClient, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	}) 

        pages, err := volumes.List(volumeClient, volumes.ListOpts{
            Name: name,
            Metadata: map[string]string{
                "migrate_kit": "true",
                "vm":          vm,
                "disk":        disk,
            },
	    }).AllPages(ctx)

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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
