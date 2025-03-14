package storage

// Data Access Object Storage
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: Update the New function to implement the creation of a new domain entity
//TODO: Create any new functions required to support the domain entity

import (
	"context"
	"fmt"
	"time"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func New(ctx context.Context, inName, inMountPoint, inDevice, inType, inOptions, inHost, inIP string) (Storage_Store, error) {

	dao.CheckDAOReadyState(domain, audit.CREATE, initialised) // Check the DAO has been initialised, Mandatory.

	//logHandler.InfoLogger.Printf("New %v (%v=%v)", domain, FIELD_ID, field1)
	clock := timing.Start(domain, actions.CREATE.GetCode(), fmt.Sprintf("%v", inMountPoint))

	// Create a new struct
	record := Storage_Store{}
	record.Raw = inHost + cfg.SEP() + inMountPoint
	record.Key = idHelpers.Encode(record.Raw)

	//m.Mountpoint, m.Source, m.FSType, Hosts
	record.MountPoint = inMountPoint
	record.Device = inDevice
	record.FSType = inType
	record.Host = inHost
	record.HostIP = inIP
	record.Options = inOptions
	record.Name = inName
	record.LastMonitored = time.Now()
	record.EverMonitored.Set(false)
	// Record the create action in the audit data
	auditErr := record.Audit.Action(ctx, audit.CREATE.WithMessage(fmt.Sprintf("New %v created %v", domain, record.Raw)))
	if auditErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOUpdateAuditError(domain, record.ID, auditErr))
	}

	// Save the status instance to the database
	writeErr := activeDB.Create(&record)
	if writeErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOCreateError(domain, record.ID, writeErr))
		//	panic(writeErr)
	}

	//logHandler.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.CREATE, domain, record.ID, fmt.Sprintf("New %v: %v", domain, field1))
	clock.Stop(1)
	return record, nil
}
