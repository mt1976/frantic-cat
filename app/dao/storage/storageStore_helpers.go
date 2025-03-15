package storage

import (
	"context"
	"fmt"
	"time"

	reporthandler "github.com/mt1976/frantic-cat/app/business/reportHandler"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

// Data Access Object Storage
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: Implement the validate function to process the domain entity
//TODO: Implement the calculate function to process the domain entity
//TODO: Implement the isDuplicateOf function to process the domain entity
//TODO: Implement the postGetProcessing function to process the domain entity

func (record *Storage_Store) upgradeProcessing() error {
	//TODO: Add any upgrade processing here
	//This processing is triggered directly after the record has been retrieved from the database
	return nil
}

func (record *Storage_Store) defaultProcessing() error {
	//TODO: Add any default calculations here
	//This processing is triggered directly before the record is saved to the database
	return nil
}

func (record *Storage_Store) validationProcessing() error {
	//TODO: Add any record validation here
	//This processing is triggered directly before the record is saved to the database and after the default calculations
	return nil
}

func (h *Storage_Store) postGetProcessing() error {
	//TODO: Add any post get processing here
	//This processing is triggered directly after the record has been retrieved from the database and after the upgrade processing
	return nil
}

func (record *Storage_Store) preDeleteProcessing() error {
	//TODO: Add any pre delete processing here
	//This processing is triggered directly before the record is deleted from the database
	return nil
}

func cloneProcessing(ctx context.Context, source Storage_Store) (Storage_Store, error) {
	//TODO: Add any clone processing here
	panic("Not Implemented")
	return Storage_Store{}, nil
}

func jobProcessor(j jobs.Job) {
	clock := timing.Start(jobs.CodedName(j), actions.PROCESS.GetCode(), j.Description())
	count := 0

	//TODO: Add your job processing code here

	// Get all the sessions
	// For each session, check the expiry date
	// If the expiry date is less than now, then delete the session

	report, err := reporthandler.NewReport(j.Name() + " Report")
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error=[%v]", jobs.CodedName(j), err.Error())
		return
	}

	StorageEntries, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error=[%v]", jobs.CodedName(j), err.Error())
		return
	}

	noStorageEntries := len(StorageEntries)
	if noStorageEntries == 0 {
		logHandler.ServiceLogger.Printf("[%v] No %vs to process", jobs.CodedName(j), domain)
		clock.Stop(0)
		return
	}

	report.AddRow(fmt.Sprintf("Found %v device(s)", noStorageEntries))

	activeEntries, err := Catalog(cfg, false)
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error=[%v]", jobs.CodedName(j), err.Error())
		return
	}

	if len(activeEntries) == 0 {
		logHandler.ServiceLogger.Printf("[%v] No %vs to process", jobs.CodedName(j), domain)
		report.AddRow("No active devices found to process")
	}

	if noStorageEntries > len(activeEntries) || noStorageEntries < len(activeEntries) {
		logHandler.ServiceLogger.Printf("[%v] %v %vs to process, but %v %vs found", jobs.CodedName(j), noStorageEntries, domain, len(activeEntries), domain)
		report.AddRow(fmt.Sprintf("%v device(s) on record, but %v active device(s) found", noStorageEntries, len(activeEntries)))
	}

	jobInstance, err := idHelpers.GetUUIDv2WithPayload(host)
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error=[%v]", jobs.CodedName(j), err.Error())
		return
	}

	report.H1("Checking Active Devices")

	for StorageEntryIndex, StorageRecord := range StorageEntries {
		logHandler.ServiceLogger.Printf("[%v] %v(%v/%v) %v", jobs.CodedName(j), domain, StorageEntryIndex+1, noStorageEntries, StorageRecord.Raw)
		StorageRecord.Signature = jobInstance
		StorageRecord.LastMonitored = time.Now()
		StorageRecord.EverMonitored.Set(true)

		// Check that this entry is in the list of active entries
		// If it is not, then log it

		//	fmt.Printf("activeEntries: %+v\n", activeEntries)

		found := find(StorageRecord, activeEntries)
		if !found {
			logHandler.WarningLogger.Printf("[%v] %v(%v/%v) %v not found in active entries", jobs.CodedName(j), domain, StorageEntryIndex+1, noStorageEntries, StorageRecord.Raw)
			report.AddRow(fmt.Sprintf("'%v' not found in active devices (%v)", StorageRecord.Name, StorageRecord.MountPoint))
			// send a notification
		} else {
			err := StorageRecord.UpdateWithAction(context.TODO(), audit.SILENT, "")
			if err != nil {
				logHandler.ErrorLogger.Printf("[%v] Error=[%v]", jobs.CodedName(j), err.Error())
				continue
			}
			report.AddRow(fmt.Sprintf("'%v' found in active devices (%v)", StorageRecord.Name, StorageRecord.MountPoint))
		}

		// check that the current item is still active
		//StorageRecord.UpdateWithAction(context.TODO(), audit.GRANT, "Job Processing")
		//StorageRecord.UpdateWithAction(context.TODO(), audit.SERVICE, "Job Processing "+j.Name())
		//count++
	}
	_ = report.Spool()
	clock.Stop(count)
}

func find(record Storage_Store, list []Storage_Store) bool {
	for _, item := range list {
		//logHandler.ServiceLogger.Printf("Comparing %v with %v", record.Raw, item.Raw)
		if record.Raw == item.Raw {
			logHandler.ServiceLogger.Printf("Found %v in %v", record.Raw, item.Raw)
			return true
		}
	}
	logHandler.WarningLogger.Printf("Did not find %v in list", record.Raw)
	return false
}
