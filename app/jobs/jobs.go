package jobs

import (
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
)

// var TemplateJobInstance jobs.Job = &TemplateJob{}
var StorageMonitorJobInstance jobs.Job = &StorageMonitorJob{}
var ProbeJobInstance jobs.Job = &ProbeJob{}

func Start() {
	logHandler.ServiceLogger.Println("Starting Jobs")
	err := jobs.Initialise(cfg)
	if err != nil {
		panic(err)
	}
	jobs.AddJobToScheduler(StorageMonitorJobInstance)

	jobs.AddJobToScheduler(ProbeJobInstance)

	jobs.StartScheduler()

	logHandler.SecurityLogger.Println("Jobs started")
}
