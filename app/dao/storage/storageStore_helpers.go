package storage

import (
	"context"

	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
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

	for StorageEntryIndex, StorageRecord := range StorageEntries {
		logHandler.ServiceLogger.Printf("[%v] %v(%v/%v) %v", jobs.CodedName(j), domain, StorageEntryIndex+1, noStorageEntries, StorageRecord.Raw)
		StorageRecord.UpdateWithAction(context.TODO(), audit.GRANT, "Job Processing")
		StorageRecord.UpdateWithAction(context.TODO(), audit.SERVICE, "Job Processing "+j.Name())
		count++
	}
	clock.Stop(count)
}
