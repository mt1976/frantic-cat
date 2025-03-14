package storage

// Data Access Object Storage
// Version: 0.2.0
// Updated on: 2021-09-10

//TODO: Update the Storage_Store struct to match the domain entity
//TODO: Update the FIELD_ constants to match the domain entity

import (
	"time"

	"github.com/mt1976/frantic-core/dao"
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
	MountPoint    string        ``              // mount point
	Device        string        ``              // source
	FSType        string        ``              // file system type
	Options       string        `csv:"-"`       // options
	Host          string        `storm:"index"` // origin
	HostIP        string        `storm:"index"` // origin IP
	Name          string        ``              // name
	Signature     string        `csv:"-"`       // signature
	LastMonitored time.Time     `csv:"-"`       // last seen
	EverMonitored dao.StormBool `csv:"-"`       // ever monitored
	//m.Mountpoint, m.Source, m.FSType

}

// Define the field set as names
var (
	FIELD_ID    = "ID"
	FIELD_Key   = "Key"
	FIELD_Raw   = "Raw"
	FIELD_Audit = "Audit"
	// Add your fields here
	FIELD_MountPoint    = "MountPoint"
	FIELD_Device        = "Device"
	FIELD_FSType        = "FSType"
	FIELD_Options       = "Options"
	FIELD_Host          = "Host"
	FIELD_HostIP        = "HostIP"
	FIELD_Name          = "Name"
	FIELD_Signature     = "Signature"
	FIELD_LastMonitored = "LastMonitored"
	FIELD_EverMonitored = "EverMonitored"
)
