package probeHandler

import (
	"fmt"

	reporthandler "github.com/mt1976/frantic-cat/app/business/reportHandler"
	"github.com/mt1976/frantic-core/application"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/netHandler"
	"github.com/mt1976/frantic-core/notificationHandler"
	"github.com/mt1976/frantic-core/timing"
)

// Worker is a job that is scheduled to run at a predefined interval
func Worker(j jobs.Job, db *database.DB) {
	clock := timing.Start(jobs.CodedName(j), actions.INITIALISE.GetCode(), j.Description())

	jobProcessor(j)

	clock.Stop(1)
}

func jobProcessor(j jobs.Job) {

	hostname := application.HostName()

	report, err := reporthandler.NewReport("Host Availability", reporthandler.TYPE_Default)
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
	var msg string
	var title string

	report.Break()

	switch {
	case noHosts == count:
		logHandler.ServiceLogger.Println("All hosts are up")
		report.AddRow("All hosts are up")
		msg = fmt.Sprintf("%v can connect to external sites", hostname)
		title = fmt.Sprintf("%v - has external connections", hostname)
	case count == 0:
		logHandler.WarningLogger.Println("All hosts are down")
		report.AddRow("All hosts are down")
		msg = fmt.Sprintf("%v cannot connect to external sites", hostname)
		title = fmt.Sprintf("%v - no external connections", hostname)
	default:
		msg = fmt.Sprintf("%v cannot connect some sites", hostname)
		title = fmt.Sprintf("%v - some hosts are inaccesible", hostname)
		logHandler.WarningLogger.Printf("Some hosts are down, %v of %v hosts down!", noHosts-count, noHosts)
		report.AddRow(fmt.Sprintf("Some hosts are down, %v of %v hosts down!", noHosts-count, noHosts))
	}

	_ = notificationHandler.Send(msg, title, 0)

	report.Spool()
	clock.Stop(count)
}
