package jobs

import (
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
)

// var TemplateJobInstance jobs.Job = &TemplateJob{}
var StorageMonitorJobInstance jobs.Job = &StorageMonitorJob{}
var ProbeJobInstance jobs.Job = &ProbeJob{}

func Start() {

	// Run the jobs once

	err := StorageMonitorJobInstance.Run()
	if err != nil {
		logHandler.ServiceLogger.Println("Error in StorageMonitorJobInstance: ", err)
	}

	err = ProbeJobInstance.Run()
	if err != nil {
		logHandler.ServiceLogger.Println("Error in ProbeJobInstance: ", err)
	}

	logHandler.ServiceLogger.Println("Starting Jobs")
	err = jobs.Initialise(cfg)
	if err != nil {
		panic(err)
	}
	jobs.AddJobToScheduler(StorageMonitorJobInstance)

	jobs.AddJobToScheduler(ProbeJobInstance)

	jobs.StartScheduler()

	logHandler.ServiceLogger.Println("Jobs started")
}
