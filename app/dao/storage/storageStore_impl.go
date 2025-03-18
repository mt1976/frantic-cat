package storage

import (
	"context"
	"fmt"
	"os"
	"strings"

	reporthandler "github.com/mt1976/frantic-cat/app/business/reportHandler"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/shirou/gopsutil/disk"
)

func Catalog(cfg *commonConfig.Settings, catalogData bool) ([]Storage_Store, error) {
	// This is the main function
	if catalogData {
		err := ClearDown(context.TODO())
		if err != nil {
			logHandler.ErrorLogger.Println("Error dropping storage records: ", err)
			panic(err)
		}
	}

	// host := application.HostName()
	// hostIP := application.HostIP()
	if catalogData {
		logHandler.InfoLogger.Println("Running in Catalog Mode")
	} else {
		logHandler.InfoLogger.Println("Running in Job Mode")
	}

	report, _ := reporthandler.NewReport("Storage Catalog", reporthandler.TYPE_Default)

	disks, err := disk.Partitions(true)
	if err != nil {
		logHandler.ErrorLogger.Println("Error getting disks: ", err)
		report.AddRow("Error getting disks: " + err.Error())
		report.Spool()
		panic(err)
	}

	var thrombuses []Storage_Store

	report.AddRow(fmt.Sprintf("Number of disks: %v (Includes TimeMachine mounts)", len(disks)))
	report.H1("Storage Catalog")

	for _, m := range disks {

		if m.Fstype == "nullfs" || m.Fstype == "overlay" {
			logHandler.InfoLogger.Printf("Skipping %v mount: %v", m.Fstype, m.Mountpoint)
			continue
		}
		// Should we skip this mount? yes, if the mountpoint contains ".timemachine"
		if strings.Contains(m.Mountpoint, ".timemachine") || strings.Contains(m.Mountpoint, ".TimeMachine") {
			logHandler.InfoLogger.Printf("Skipping TimeMachine mount: %v", m.Mountpoint)
			continue
		}

		//	fmt.Printf("Mount=%v Source=%v Type=%v\n", m.Mountpoint, m.Source, m.FSType)
		//logHandler.InfoLogger.Printf("Data=%+v\n", m)
		name := getMountName(m.Mountpoint)
		if catalogData {
			logHandler.EventLogger.Printf("Adding mount record: %v '%v'", name, m.Mountpoint)

			thrombus, err := New(context.TODO(), name, m.Mountpoint, m.Device, m.Fstype, m.Opts, host, hostIP)
			if err != nil {
				logHandler.ErrorLogger.Println("Error creating storage record: ", err)
				panic(err)
			}
			thrombuses = append(thrombuses, thrombus)
		} else {
			thrombus := Storage_Store{}
			thrombus.Name = name
			thrombus.Raw = host + cfg.SEP() + m.Mountpoint
			thrombus.Key = idHelpers.Encode(thrombus.Raw)
			thrombus.MountPoint = m.Mountpoint
			thrombus.Device = m.Device
			thrombus.FSType = m.Fstype
			thrombus.Host = host
			thrombus.HostIP = hostIP
			thrombus.Options = m.Opts
			thrombuses = append(thrombuses, thrombus)

		}
		report.AddRow(fmt.Sprintf("Mount: %v", m.Mountpoint))
	}
	if catalogData {
		report.HR()
		report.AddRow(fmt.Sprintf("Number of mounts: %v", len(thrombuses)))
		report.HR()
		err = ExportRecordsAsCSV()

		if err != nil {
			logHandler.ErrorLogger.Println("Error exporting storage records: ", err)
			panic(err)
		}
		report.Spool()
	}
	return thrombuses, nil
}

func getMountName(m string) string {
	name := m
	//get last element of name delimited by os.PathSeparator
	names := strings.Split(name, string(os.PathSeparator))
	name = names[len(names)-1:][0]
	if len(name) == 0 {
		name = m
	}
	return name
}
