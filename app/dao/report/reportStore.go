package report

// Data Access Object Report
// Version: 0.2.0
// Updated on: 2021-09-10

/// DO NOT CHANGE THIS FILE OTHER THAN TO RENAME "Report" TO THE NAME OF THE DOMAIN ENTITY
/// DO NOT CHANGE THIS FILE OTHER THAN TO RENAME "Report" TO THE NAME OF THE DOMAIN ENTITY
/// DO NOT CHANGE THIS FILE OTHER THAN TO RENAME "Report" TO THE NAME OF THE DOMAIN ENTITY

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/dao/lookup"
	"github.com/mt1976/frantic-core/importExportHelper"
	"github.com/mt1976/frantic-core/ioHelpers"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/paths"
	"github.com/mt1976/frantic-core/timing"
)

func Count() (int, error) {
	logHandler.DatabaseLogger.Printf("Count %v", domain)
	return activeDB.Count(&Report_Store{})
}

func CountWhere(field string, value any) (int, error) {
	logHandler.DatabaseLogger.Printf("Count %v where (%v=%v)", domain, field, value)
	clock := timing.Start(domain, actions.COUNT.GetCode(), fmt.Sprintf("%v=%v", field, value))
	list, err := GetAllWhere(field, value)
	if err != nil {
		clock.Stop(0)
		return 0, err
	}
	clock.Stop(len(list))
	return len(list), nil
}

func GetById(id any) (Report_Store, error) {
	return GetBy(FIELD_ID, id)
}

func GetByKey(key any) (Report_Store, error) {
	return GetBy(FIELD_Key, key)
}

func GetBy(field string, value any) (Report_Store, error) {

	clock := timing.Start(domain, actions.GET.GetCode(), fmt.Sprintf("%v=%v", field, value))

	dao.CheckDAOReadyState(domain, audit.GET, initialised) // Check the DAO has been initialised, Mandatory.

	if field == FIELD_ID && reflect.TypeOf(value).Name() != "int" {
		msg := "invalid data type. Expected type of %v is int"
		logHandler.ErrorLogger.Printf(msg, value)
		return Report_Store{}, commonErrors.WrapDAOReadError(domain, field, value, fmt.Errorf(msg, value))
	}

	if err := dao.IsValidFieldInStruct(field, Report_Store{}); err != nil {
		return Report_Store{}, err
	}

	if err := dao.IsValidTypeForField(field, value, Report_Store{}); err != nil {
		return Report_Store{}, err
	}

	record := Report_Store{}
	logHandler.DatabaseLogger.Printf("Get %v where (%v=%v)", domain, field, value)

	if err := activeDB.Retrieve(field, value, &record); err != nil {
		clock.Stop(0)
		return Report_Store{}, commonErrors.WrapRecordNotFoundError(domain, field, fmt.Sprintf("%v", value))
	}

	if err := record.postGet(); err != nil {
		clock.Stop(0)
		return Report_Store{}, commonErrors.WrapDAOReadError(domain, field, value, err)
	}

	clock.Stop(1)
	return record, nil
}

func GetAll() ([]Report_Store, error) {

	dao.CheckDAOReadyState(domain, audit.GET, initialised) // Check the DAO has been initialised, Mandatory.

	recordList := []Report_Store{}

	clock := timing.Start(domain, actions.GETALL.GetCode(), "ALL")

	if errG := activeDB.GetAll(&recordList); errG != nil {
		clock.Stop(0)
		return []Report_Store{}, commonErrors.WrapNotFoundError(domain, errG)
	}

	var errPost error
	if recordList, errPost = postGetList(&recordList); errPost != nil {
		clock.Stop(0)
		return nil, errPost
	}

	clock.Stop(len(recordList))

	return recordList, nil
}

func GetAllWhere(field string, value any) ([]Report_Store, error) {
	dao.CheckDAOReadyState(domain, audit.GET, initialised) // Check the DAO has been initialised, Mandatory.

	recordList := []Report_Store{}
	resultList := []Report_Store{}

	clock := timing.Start(domain, actions.GETALL.GetCode(), fmt.Sprintf("%v=%v", field, value))

	if err := dao.IsValidFieldInStruct(field, Report_Store{}); err != nil {
		return recordList, err
	}

	if err := dao.IsValidTypeForField(field, value, Report_Store{}); err != nil {
		return recordList, err
	}

	//err := activeDB.Retrieve(field, value, &recordList)

	recordList, err := GetAll()
	if err != nil {
		return []Report_Store{}, err
	}
	count := 0

	for _, record := range recordList {
		if reflect.ValueOf(record).FieldByName(field).Interface() == value {
			count++
			resultList = append(resultList, record)
		}
	}

	var errPost error
	if resultList, errPost = postGetList(&resultList); errPost != nil {
		clock.Stop(0)
		return nil, errPost
	}

	clock.Stop(len(resultList))

	return resultList, nil
}

func Delete(ctx context.Context, id int, note string) error {
	return DeleteBy(ctx, FIELD_ID, id, note)
}

func DeleteByKey(ctx context.Context, key string, note string) error {
	return DeleteBy(ctx, FIELD_Key, key, note)
}
func DeleteBy(ctx context.Context, field string, value any, note string) error {

	dao.CheckDAOReadyState(domain, audit.DELETE, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.DELETE.GetCode(), fmt.Sprintf("%v=%v", field, value))

	if err := dao.IsValidFieldInStruct(field, Report_Store{}); err != nil {
		logHandler.ErrorLogger.Print(commonErrors.WrapDAODeleteError(domain, field, value, err).Error())
		clock.Stop(0)
		return commonErrors.WrapDAODeleteError(domain, field, value, err)
	}

	if err := dao.IsValidTypeForField(field, value, Report_Store{}); err != nil {
		logHandler.ErrorLogger.Print(commonErrors.WrapDAODeleteError(domain, field, value, err).Error())
		clock.Stop(0)
		return err
	}

	record, err := GetBy(field, value)

	if err != nil {
		getErr := commonErrors.WrapDAODeleteError(domain, field, value, err)
		logHandler.ErrorLogger.Panic(getErr.Error(), err)
		clock.Stop(0)
		return getErr
	}

	auditErr := record.Audit.Action(ctx, audit.DELETE.WithMessage(note))
	if auditErr != nil {
		audErr := commonErrors.WrapDAOUpdateAuditError(domain, value, auditErr)
		logHandler.ErrorLogger.Print(audErr.Error())
		clock.Stop(0)
		return audErr
	}

	preDeleteErr := record.preDeleteProcessing()
	if preDeleteErr != nil {
		logHandler.ErrorLogger.Print(commonErrors.WrapDAODeleteError(domain, field, value, preDeleteErr).Error())
		clock.Stop(0)
		return preDeleteErr
	}

	record.ExportRecordAsJSON(audit.DELETE.Description())

	if err := activeDB.Delete(&record); err != nil {
		delErr := commonErrors.WrapDAODeleteError(domain, field, value, err)
		logHandler.ErrorLogger.Panic(delErr.Error())
		clock.Stop(0)
		return delErr
	}

	clock.Stop(1)

	return nil
}

func (record *Report_Store) Spew() {
	logHandler.InfoLogger.Printf("[%v] Record=[%+v]", domain, record)
}

func (record *Report_Store) Validate() error {
	return record.validationProcessing()
}

func (record *Report_Store) Update(ctx context.Context, note string) error {
	return record.insertOrUpdate(ctx, note, actions.UPDATE.GetCode(), audit.UPDATE, actions.UPDATE.GetCode())
}

func (record *Report_Store) UpdateWithAction(ctx context.Context, auditAction audit.Action, note string) error {
	return record.insertOrUpdate(ctx, note, actions.UPDATE.GetCode(), auditAction, actions.UPDATE.GetCode())
}

func (record *Report_Store) Create(ctx context.Context, note string) error {
	return record.insertOrUpdate(ctx, note, actions.CREATE.GetCode(), audit.CREATE, actions.CREATE.GetCode())
}

func (record *Report_Store) Clone(ctx context.Context) (Report_Store, error) {
	return reportCloneProcessing(ctx, *record)
}

func (record *Report_Store) ExportRecordAsJSON(name string) {

	ID := reflect.ValueOf(*record).FieldByName(FIELD_ID)

	clock := timing.Start(domain, actions.EXPORT.GetCode(), fmt.Sprintf("%v", ID))

	ioHelpers.Dump(domain, paths.Dumps(), name, fmt.Sprintf("%v", ID), record)

	clock.Stop(1)
}

func GetDefaultLookup() (lookup.Lookup, error) {
	return GetLookup(FIELD_Key, FIELD_Raw)
}

func GetLookup(field, value string) (lookup.Lookup, error) {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.LOOKUP.GetCode(), "BUILD")

	// Get all status
	recordList, err := GetAll()
	if err != nil {
		lkpErr := commonErrors.WrapDAOLookupError(domain, field, value, err)
		logHandler.ErrorLogger.Print(lkpErr.Error())
		clock.Stop(0)
		return lookup.Lookup{}, lkpErr
	}

	// Create a new Lookup
	var rtnLookup lookup.Lookup
	rtnLookup.Data = make([]lookup.LookupData, 0)

	// range through Behaviour list, if status code is found and deletedby is empty then return error
	for _, a := range recordList {
		key := reflect.ValueOf(a).FieldByName(field).Interface().(string)
		val := reflect.ValueOf(a).FieldByName(value).Interface().(string)
		rtnLookup.Data = append(rtnLookup.Data, lookup.LookupData{Key: key, Value: val})
	}

	clock.Stop(len(rtnLookup.Data))

	return rtnLookup, nil
}

func Drop() error {
	return activeDB.Drop(Report_Store{})
}

// GetDatabaseConnections returns a function that fetches the current database instances.
//
// This function is used to retrieve the active database instances being used by the application.
// It returns a function that, when called, returns a slice of pointers to `database.DB` and an error.
//
// Returns:
//
//	func() ([]*database.DB, error): A function that returns a slice of pointers to `database.DB` and an error.
func GetDatabaseConnections() func() ([]*database.DB, error) {
	return func() ([]*database.DB, error) {
		return []*database.DB{activeDB}, nil
	}
}

func ClearDown(ctx context.Context) error {
	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.CLEAR.GetCode(), "INITIALISE")

	// Delete all active session recordList
	recordList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Print(commonErrors.WrapDAOInitialisationError(domain, err).Error())
		clock.Stop(0)
		return commonErrors.WrapDAOInitialisationError(domain, err)
	}

	noRecords := len(recordList)
	count := 0

	for thisRecord, record := range recordList {
		logHandler.InfoLogger.Printf("Deleting %v (%v/%v) %v", domain, thisRecord, noRecords, record.Key)
		delErr := Delete(ctx, record.ID, fmt.Sprintf("Clearing %v %v @ initialisation ", domain, record.ID))
		if delErr != nil {
			logHandler.ErrorLogger.Print(commonErrors.WrapDAOInitialisationError(domain, delErr).Error())
			continue
		}
		count++
	}

	clock.Stop(count)
	logHandler.EventLogger.Printf("Cleared down %v", domain)
	return nil
}

func ExportRecordsAsJSON(message string) {

	dao.CheckDAOReadyState(domain, audit.EXPORT, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.EXPORT.GetCode(), "ALL")
	recordList, _ := GetAll()
	if len(recordList) == 0 {
		logHandler.WarningLogger.Printf("[%v] %v data not found", domain, domain)
		clock.Stop(0)
		return
	}
	SEP := "!"
	for _, record := range recordList {
		msg := fmt.Sprintf("%v%v%v", audit.EXPORT.Description(), SEP, message)
		if message == "" {
			msg = fmt.Sprintf("%v%v", audit.EXPORT.Description(), SEP)
		}
		record.ExportRecordAsJSON(msg)
	}
	clock.Stop(len(recordList))
}

func ExportRecordsAsCSV() error {

	exportListData, err := GetAll()
	if err != nil {
		logHandler.ExportLogger.Panicf("error Getting all %v's: %v", domain, err.Error())
	}

	return importExportHelper.ExportCSV(domain, exportListData)
}

func ImportRecordsFromCSV() error {
	return importExportHelper.ImportCSV(domain, &Report_Store{}, importProcessor)
}

// Worker is a job that is scheduled to run at a predefined interval
func Worker(j jobs.Job, db *database.DB) {
	clock := timing.Start(jobs.CodedName(j), actions.INITIALISE.GetCode(), j.Description())
	oldDB := activeDB
	dbSwitched := false
	// Overide the default database connection if one is passed

	if db != nil {
		if activeDB.Name != db.Name {
			logHandler.EventLogger.Printf("Switching to %v.db", db.Name)
			activeDB = db
			dbSwitched = true
		}
	}

	reportJobProcessor(j)

	if dbSwitched {
		logHandler.EventLogger.Printf("Switching back to %v.db from %v.db", oldDB.Name, activeDB.Name)
		activeDB = oldDB
	}
	clock.Stop(1)
}
