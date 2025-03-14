package storage

// Data Access Object Storage
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"context"
	"fmt"
	"strings"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func (record *Storage_Store) insertOrUpdate(ctx context.Context, note, activity string, auditAction audit.Action, operation string) error {

	isCreateOperation := false
	if strings.EqualFold(operation, actions.CREATE.GetCode()) {
		isCreateOperation = true
		if !strings.EqualFold(auditAction.Code(), actions.CREATE.GetCode()) {
			return commonErrors.WrapDAOUpdateError(domain, fmt.Errorf("invalid audit action '%v' for create event '%v'", auditAction.Code(), operation))
		}
	}

	dao.CheckDAOReadyState(domain, auditAction, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, activity, fmt.Sprintf("%v", record.ID))
	if isCreateOperation {
		if err := record.checkForDuplicate(); err != nil {
			clock.Stop(0)
			return commonErrors.WrapDAOCreateError(domain, record.ID, err)
		}
	}

	if calculationError := record.defaultProcessing(); calculationError != nil {
		rtnErr := commonErrors.WrapDAOCaclulationError(domain, calculationError)
		logHandler.ErrorLogger.Print(rtnErr.Error())
		clock.Stop(0)
		return rtnErr
	}

	if validationError := record.validationProcessing(); validationError != nil {
		valErr := commonErrors.WrapDAOValidationError(domain, validationError)
		logHandler.ErrorLogger.Print(valErr.Error())
		clock.Stop(0)
		return valErr
	}

	auditErr := record.Audit.Action(ctx, auditAction.WithMessage(note))
	if auditErr != nil {
		audErr := commonErrors.WrapDAOUpdateAuditError(domain, record.ID, auditErr)
		logHandler.ErrorLogger.Print(audErr.Error())
		clock.Stop(0)
		return audErr
	}
	var actionError error
	if isCreateOperation {
		actionError = activeDB.Create(record)
	} else {
		actionError = activeDB.Update(record)
	}
	if actionError != nil {
		updErr := commonErrors.WrapDAOUpdateError(domain, actionError)
		logHandler.ErrorLogger.Panic(updErr.Error(), actionError)
		clock.Stop(0)
		return updErr
	}

	clock.Stop(1)

	return nil
}

func postGetList(recordList *[]Storage_Store) ([]Storage_Store, error) {
	clock := timing.Start(domain, actions.PROCESS.GetCode(), "POSTGET")
	returnList := []Storage_Store{}
	for _, record := range *recordList {
		if err := record.postGet(); err != nil {
			return nil, err
		}
		returnList = append(returnList, record)
	}
	clock.Stop(len(returnList))
	return returnList, nil
}

func (record *Storage_Store) postGet() error {

	upgradeError := record.upgradeProcessing()
	if upgradeError != nil {
		return upgradeError
	}

	defaultingError := record.defaultProcessing()
	if defaultingError != nil {
		return defaultingError
	}

	validationError := record.validationProcessing()
	if validationError != nil {
		return validationError
	}

	return record.postGetProcessing()
}

func (record *Storage_Store) checkForDuplicate() error {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	// Get all status
	responseRecord, err := GetBy(FIELD_Key, record.Key)
	if err != nil {
		// This is ok, no record could be read
		return nil
	}

	if responseRecord.Audit.DeletedBy != "" {
		return nil
	}

	logHandler.WarningLogger.Printf("Duplicate %v, %v already in use", domain, record.ID)
	return commonErrors.ErrorDuplicate
}
