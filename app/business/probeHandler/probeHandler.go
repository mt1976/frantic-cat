package probeHandler

import (
	"fmt"

	reporthandler "github.com/mt1976/frantic-cat/app/business/reportHandler"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/netHandler"
	"github.com/mt1976/frantic-core/timing"
)

// Worker is a job that is scheduled to run at a predefined interval
func Worker(j jobs.Job, db *database.DB) {
	clock := timing.Start(jobs.CodedName(j), actions.INITIALISE.GetCode(), j.Description())

	jobProcessor(j)

	clock.Stop(1)
}

func jobProcessor(j jobs.Job) {

	// TODO: this should get a list of hosts to test from the cfg.Settings and test them all
	// for now we are just testing localhost

	report, err := reporthandler.NewReport("Host Availability")
	if err != nil {
		panic(err)
	}

	cfg := commonConfig.Get()
	hosts := cfg.GetValidHosts()

	clock := timing.Start(jobs.CodedName(j), actions.PROCESS.GetCode(), j.Description())
	count := 0

	noHosts := len(hosts)
	report.AddRow(fmt.Sprintf("Number of hosts to test: %v", noHosts))
	report.H1("Host Availability")
	for _, host := range hosts {
		//	up, err := minPing(host.FQDN, application.IsRunningOnWindows())
		up, err := netHandler.CheckHostAvailability(host.FQDN)

		if err != nil {
			continue
		}
		if up {
			report.AddRow(fmt.Sprintf("%v is up", host.FQDN))
			count++
		} else {
			report.AddRow(fmt.Sprintf("%v is down", host.FQDN))
		}
	}
	report.Break()
	switch {
	case noHosts == count:
		fmt.Println("All hosts are up")
		report.AddRow("All hosts are up")
	case count == 0:
		fmt.Println("All hosts are down")
		report.AddRow("All hosts are down")
	default:
		fmt.Println("Some hosts are down")
		report.AddRow("Some hosts are down")
		report.AddRow(fmt.Sprintf("%v of %v hosts down!", noHosts-count, noHosts))
	}
	report.Spool()
	clock.Stop(count)
}
