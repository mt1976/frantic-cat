package main

import (
	"context"

	reportStore "github.com/mt1976/frantic-cat/app/dao/report"
	storageStore "github.com/mt1976/frantic-cat/app/dao/storage"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/spf13/pflag"
)

func main() {
	// This is the main function

	// Going to need to run this with command line arguments
	// 1. --catalog - generates the initial catalog (or re-generates it)
	// 2. No arguments - runs the job

	var inCatalogMode *bool = pflag.Bool("catalog", false, "Generate the initial catalog")
	pflag.Parse()

	if *inCatalogMode {
		logHandler.InfoLogger.Println("Running in Catalog Mode")
	} else {
		logHandler.InfoLogger.Println("Running in Job Mode")
	}

	storageStore.Initialise(context.TODO())
	reportStore.Initialise(context.TODO())

	cfg := commonConfig.Get()

	if *inCatalogMode {
		err := Catalog(cfg, true)
		if err != nil {
			logHandler.ErrorLogger.Println("Error in Catalog Mode: ", err)
			panic(err)
		}

	} else {
		Monitor(cfg)
	}

	logHandler.InfoLogger.Println("Finished")

}
