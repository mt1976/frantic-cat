package template

// Data Access Object Template
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: RENAME "Template" TO THE NAME OF THE DOMAIN ENTITY
//TODO: Update the Template_Store struct to match the domain entity
//TODO: Update the FIELD_ constants to match the domain entity

import (
	"time"

	audit "github.com/mt1976/frantic-core/dao/audit"
)

var domain = "Template"

// Template_Store represents a Template_Store entity.
type Template_Store struct {
	ID    int         `storm:"id,increment=100000"` // primary key with auto increment
	Key   string      `storm:"unique"`              // key
	Raw   string      `storm:"unique"`              // raw ID before encoding
	Audit audit.Audit `csv:"-"`                     // audit data
	// Add your fields here
	Field1 int       `csv:"altID" storm:"index"` // user key
	Field2 string    `storm:"index"`             // user code
	Field3 time.Time `csv:"-"`                   // expiry time
}

// Define the field set as names
var (
	FIELD_ID    = "ID"
	FIELD_Key   = "Key"
	FIELD_Raw   = "Raw"
	FIELD_Audit = "Audit"
	// Add your fields here
	FIELD_Field1 = "Field1"
	FIELD_Field2 = "Field2"
	FIELD_Field3 = "Field3"
)
