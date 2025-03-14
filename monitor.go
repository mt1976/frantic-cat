package main

import (
	"github.com/mt1976/frantic-cat/app/jobs"
	"github.com/mt1976/frantic-core/commonConfig"
)

func Monitor(cfg *commonConfig.Settings) {
	// This is the main function
	x := jobs.StorageMonitorJob{}
	err := x.Run()
	if err != nil {
		panic(err)
	}
}
