package main

import (
	"github.com/mt1976/frantic-cat/app/jobs"
	"github.com/mt1976/frantic-core/commonConfig"
)

func Monitor(cfg *commonConfig.Settings) {
	jobs.Start()
}
