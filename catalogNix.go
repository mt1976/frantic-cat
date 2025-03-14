package main

import (
	"context"

	"github.com/moby/sys/mountinfo"
	"github.com/mt1976/frantic-cat/app/dao/storage"
	"github.com/mt1976/frantic-core/application"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/logHandler"
)

func CatalogNix(cfg *commonConfig.Settings) {
	// This is the main function

	err := storage.ClearDown(context.TODO())
	if err != nil {
		logHandler.ErrorLogger.Println("Error dropping storage records: ", err)
		panic(err)
	}

	host := application.HostName()
	hostIP := application.HostIP()

	mounts, err := mountinfo.GetMounts(nil)
	if err != nil {
		logHandler.ErrorLogger.Println("Error getting mounts: ", err)
		panic(err)
	}

	for _, m := range mounts {

		if m.FSType == "nullfs" || m.FSType == "overlay" {
			logHandler.InfoLogger.Printf("Skipping %v mount: %v", m.FSType, m.Mountpoint)
			continue
		}

		//	fmt.Printf("Mount=%v Source=%v Type=%v\n", m.Mountpoint, m.Source, m.FSType)

		logHandler.EventLogger.Printf("Creating mount record: %v", m.Mountpoint)
		_, err := storage.New(context.TODO(), m.Mountpoint, m.Source, m.FSType, host, hostIP)
		if err != nil {
			logHandler.ErrorLogger.Println("Error creating storage record: ", err)
			panic(err)
		}
	}

	err = storage.ExportRecordsAsCSV()

	if err != nil {
		logHandler.ErrorLogger.Println("Error exporting storage records: ", err)
		panic(err)
	}
}
