package main

import (
	"context"

	"os"
	"strings"

	"github.com/mt1976/frantic-cat/app/dao/storage"
	"github.com/mt1976/frantic-core/application"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/logHandler"
	disk "github.com/shirou/gopsutil/disk"
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

	logHandler.InfoLogger.Println("Running in Catalog Mode")

	disks, err := disk.Partitions(true)
	if err != nil {
		logHandler.ErrorLogger.Println("Error getting disks: ", err)
		panic(err)
	}

	for _, m := range disks {

		if m.Fstype == "nullfs" || m.Fstype == "overlay" {
			logHandler.InfoLogger.Printf("Skipping %v mount: %v", m.Fstype, m.Mountpoint)
			continue
		}

		//	fmt.Printf("Mount=%v Source=%v Type=%v\n", m.Mountpoint, m.Source, m.FSType)
		logHandler.InfoLogger.Printf("Data=%+v\n", m)
		name := getMountName(m.Mountpoint)
		logHandler.EventLogger.Printf("Creating mount record: %v '%v'", name, m.Mountpoint)

		_, err := storage.New(context.TODO(), name, m.Mountpoint, m.Device, m.Fstype, m.Opts, host, hostIP)
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

func getMountName(m string) string {
	name := m
	//get last element of name delimited by os.PathSeparator
	names := strings.Split(name, string(os.PathSeparator))
	name = names[len(names)-1:][0]
	return name
}
