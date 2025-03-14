package storage

// Data Access Object Storage
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: Update the Storage_Store struct to match the domain entity
//TODO: Update the FIELD_ constants to match the domain entity

import (
	audit "github.com/mt1976/frantic-core/dao/audit"
)

var domain = "Storage"

// Storage_Store represents a Storage_Store entity.
type Storage_Store struct {
	ID    int         `storm:"id,increment=100000"` // primary key with auto increment
	Key   string      `storm:"unique"`              // key
	Raw   string      `storm:"unique"`              // raw ID before encoding
	Audit audit.Audit `csv:"-"`                     // audit data
	// Add your fields here
	MountPoint string ``              // mount point
	Source     string ``              // source
	FSType     string ``              // file system type
	Host       string `storm:"index"` // origin
	HostIP     string `storm:"index"` // origin IP
	//m.Mountpoint, m.Source, m.FSType

}

// Define the field set as names
var (
	FIELD_ID    = "ID"
	FIELD_Key   = "Key"
	FIELD_Raw   = "Raw"
	FIELD_Audit = "Audit"
	// Add your fields here
	FIELD_MountPoint = "MountPoint"
	FIELD_Source     = "Source"
	FIELD_FSType     = "FSType"
	FIELD_Host       = "Host"
	FIELD_HostIP     = "HostIP"
)
