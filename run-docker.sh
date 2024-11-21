#!/bin/bash

docker run -it --rm --privileged --network host --env-file <(env | grep OS_) dorfpinguin/gophercloud-volume-list-test:latest --debug --name Testsnapsot-volume2 --vm vm-2344 --disk 1234-1000
