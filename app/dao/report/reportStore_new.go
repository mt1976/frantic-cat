package report

// Data Access Object Report
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: RENAME "Report" TO THE NAME OF THE DOMAIN ENTITY
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

func New(ctx context.Context, field1 string, field2 string) (Report_Store, error) {

	dao.CheckDAOReadyState(domain, audit.CREATE, initialised) // Check the DAO has been initialised, Mandatory.

	//logHandler.InfoLogger.Printf("New %v (%v=%v)", domain, FIELD_ID, field1)
	clock := timing.Start(domain, actions.CREATE.GetCode(), fmt.Sprintf("%v", field1))

	sessionID := idHelpers.GetUUID()

	var testTest []Row
	testTest = append(testTest, Row{Index: 1, Text: "Test1"})
	testTest = append(testTest, Row{Index: 2, Text: "Test2"})
	testTest = append(testTest, Row{Index: 3, Text: "Test3"})
	testTest = append(testTest, Row{Index: 4, Text: "Test4"})
	testTest = append(testTest, Row{Index: 5, Text: "Test5"})

	// Create a new struct
	record := Report_Store{}
	record.Key = idHelpers.Encode(sessionID)
	record.Raw = sessionID
	record.Title = field1
	record.Content = testTest

	//record.Field2 = field2
	record.Generated = time.Now().Add(time.Minute * time.Duration(sessionExpiry))

	// Record the create action in the audit data
	auditErr := record.Audit.Action(ctx, audit.CREATE.WithMessage(fmt.Sprintf("New %v created %v", domain, field1)))
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
