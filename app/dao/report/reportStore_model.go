package report

// Data Access Object Report
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: RENAME "Report" TO THE NAME OF THE DOMAIN ENTITY
//TODO: Update the Report_Store struct to match the domain entity
//TODO: Update the FIELD_ constants to match the domain entity

import (
	"time"

	audit "github.com/mt1976/frantic-core/dao/audit"
)

var domain = "Report"

// Report_Store represents a Report_Store entity.
type Report_Store struct {
	ID    int         `storm:"id,increment=100000"` // primary key with auto increment
	Key   string      `storm:"unique"`              // key
	Raw   string      `storm:"unique"`              // raw ID before encoding
	Audit audit.Audit `csv:"-"`                     // audit data
	// Add your fields here
	Title     string    `storm:"index"`  // Report Name
	Generated time.Time `storm:"index"`  // Report Generation Time
	Content   []Row     `storm:"inline"` // Report Content
	Host      string    `storm:"index"`  // Host
	HostIP    string    `storm:"index"`  // Host IP
}

type Row struct {
	Index int    // Index
	Text  string // Text
}

// Define the field set as names
var (
	FIELD_ID    = "ID"
	FIELD_Key   = "Key"
	FIELD_Raw   = "Raw"
	FIELD_Audit = "Audit"
	// Add your fields here
	FIELD_Title     = "Title"
	FIELD_Generated = "Generated"
	FIELD_Content   = "Content"
	FIELD_Index     = "Index"
	FIELD_Text      = "Text"
	FIELD_Host      = "Host"
	FIELD_HostIP    = "HostIP"
)
