package probeHandler

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mt1976/frantic-core/application"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
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

	cfg := commonConfig.Get()
	hosts := cfg.GetValidHosts()

	clock := timing.Start(jobs.CodedName(j), actions.PROCESS.GetCode(), j.Description())
	count := 0
	// up, err := miniPing("127.0.0.7", application.IsRunningOnWindows())
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// if up {
	// 	fmt.Println("IT'S ALIVEEE")
	// 	count++
	// }

	// up, err = miniPing("127.0.0.1", false)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// if up {
	// 	fmt.Println("IT'S ALIVEEE")
	// 	count++
	// }

	for _, host := range hosts {
		up, err := minPing(host.FQDN, application.IsRunningOnWindows())
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if up {
			fmt.Println("IT'S ALIVEEE")
			count++
		}
	}

	// up, err := miniPing("10.147.20.105", "7575", true)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// if up {
	// 	fmt.Println("IT'S ALIVEEE")
	// 	count++
	// }

	clock.Stop(count)
}

func minPing(addr string, windows bool) (bool, error) {
	// TODO: This should be moved to a common package
	//
	// This is a simple ping function that will return true if the host is up and false if it is down
	//
	// addr: string: The IP address of the host to ping
	// windows: bool: If the host is running windows, then we need to use a different method to ping it
	logHandler.ServiceLogger.Printf("Pinging %v", addr)
	var out []byte
	var err error
	if windows {
		logHandler.ServiceLogger.Printf("Running Windows Ping - [ping %v -n 5 -w 3000]", addr)
		out, err = exec.Command("ping", addr, "-n", "5", "-w", "3000").Output()
	} else {
		logHandler.ServiceLogger.Printf("Running Linux Ping - [ping %v -c 5 -i 3 -W 10]", addr)
		out, err = exec.Command("ping", addr, "-c 5", "-i 3", "-W 10").Output()
	}
	if err != nil {
		//		logHandler.ServiceLogger.Printf("Error: %v", err)
		if strings.Contains(err.Error(), "exit status 68") {
			//	fmt.Println("TANGO DOWN")
			logHandler.WarningLogger.Printf("Host %v is not reachable", addr)
			return false, nil
		}

	}

	if isHostReachable(out) {
		logHandler.ServiceLogger.Printf("Host %v is reachable", addr)
		//fmt.Println("TANGO DOWN")
		return true, nil
	}

	logHandler.WarningLogger.Printf("Host %v is not reachable", addr)
	return false, nil
}

func isHostReachable(out []byte) bool {

	switch {
	case strings.Contains(string(out), "Destination Host Unreachable"):
		return false
	case strings.Contains(string(out), "Request timed out"):
		return false
	case strings.Contains(string(out), "100% packet loss"):
		return false
	case strings.Contains(string(out), "Request timeout"):
		return false
	case strings.Contains(string(out), "cannot resolve"):
		return false
	}
	return true
	// b := strings.Contains(string(out), "Destination Host Unreachable") || strings.Contains(string(out), "Request timed out") || strings.Contains(string(out), "100% packet loss") || strings.Contains(string(out), "Request timed out")
	// return !b
}
