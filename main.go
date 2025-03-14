package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	templ "github.com/mt1976/frantic-cat/app/dao/template"
	"github.com/mt1976/frantic-cat/app/jobs"
	jbs "github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
)

func main() {
	// This is the main function

	templ.Initialise(context.TODO())

	// Lets clear down the session db
	initErr := templ.ClearDown(context.TODO())
	if initErr != nil {
		logHandler.ErrorLogger.Println(initErr)
	}

	logHandler.InfoLogger.Println("templateStore", "Initialise", "Done")

	iniErr := templ.ClearDown(context.TODO())
	if iniErr != nil {
		logHandler.ErrorLogger.Println(iniErr)
	}
	randLast := 0
	for i := 0; i < 10; i++ {
		randNum := rand.Intn(10000-1000) + 1000
		logHandler.InfoLogger.Println("randNum:", randNum)

		newRec, newRedErr := templ.New(context.TODO(), randNum, "test")
		if newRedErr != nil {
			logHandler.ErrorLogger.Println(newRedErr)
		} else {
			logHandler.InfoLogger.Printf("newRec:[%+v]", newRec)
		}
		newRec.ExportRecordAsJSON("test")
		randLast = randNum
	}

	templ.ExportRecordsAsJSON(time.Now().Format("20060102"))
	templ.ExportRecordsAsJSON("")

	lk, lkErr := templ.GetDefaultLookup()
	if lkErr != nil {
		logHandler.ErrorLogger.Printf("two:[%+v]", lkErr)
	} else {
		logHandler.InfoLogger.Printf("lk:[%+v]", lk)
	}

	count, cerr := templ.Count()
	if cerr != nil {
		logHandler.ErrorLogger.Println(cerr)
	} else {
		logHandler.InfoLogger.Printf("count:[%+v]", count)
	}

	count, cerr = templ.CountWhere(templ.FIELD_Field1, randLast)
	if cerr != nil {
		logHandler.ErrorLogger.Println(cerr)
	} else {
		logHandler.InfoLogger.Printf("count:[%+v]", count)
	}

	count2, cerr2 := templ.CountWhere(templ.FIELD_Field2, "test")
	if cerr2 != nil {
		logHandler.ErrorLogger.Println(cerr2)
	} else {
		logHandler.InfoLogger.Printf("count2:[%+v]", count2)
	}

	count3, cerr3 := templ.CountWhere(templ.FIELD_ID, "poopoolala")
	if cerr3 != nil {
		logHandler.ErrorLogger.Println(cerr3)
	} else {
		logHandler.InfoLogger.Printf("count3:[%+v]", count3)
	}

	count4, cerr4 := templ.CountWhere(templ.FIELD_Field1, 123)
	if cerr4 != nil {
		logHandler.ErrorLogger.Println(cerr4)
	} else {
		logHandler.InfoLogger.Printf("count4:[%+v]", count4)
	}
	/// FInal tests
	dropErr := templ.Drop()
	if dropErr != nil {
		logHandler.ErrorLogger.Println(dropErr)
	} else {
		logHandler.InfoLogger.Printf("drop:[templ]")
	}

	cdrop, cdroperr := templ.Count()
	if cdroperr != nil {
		logHandler.ErrorLogger.Println(cdroperr)
	} else {
		logHandler.InfoLogger.Printf("cdrop:[%+v]", cdrop)
	}

	for i := 0; i < 10; i++ {
		randNum := rand.Intn(10000-1000) + 1000
		//logHandler.InfoLogger.Println("randNum:", randNum)

		_, newRedErr := templ.New(context.TODO(), randNum, fmt.Sprintf("Test %v", randNum))
		if newRedErr != nil {
			logHandler.ErrorLogger.Println(newRedErr)
		}
	}

	templ.ExportRecordsAsCSV()
	// Drop Again
	dropErr = templ.ClearDown(context.TODO())
	if dropErr != nil {
		logHandler.ErrorLogger.Println(dropErr)
	} else {
		logHandler.InfoLogger.Printf("Drop Data Post Export:[templ]")
	}

	templ.ImportRecordsFromCSV()

	cdrop, cdroperr = templ.Count()
	if cdroperr != nil {
		logHandler.ErrorLogger.Println(cdroperr)
	} else {
		logHandler.InfoLogger.Printf("Post Import:[%+v]", cdrop)
	}

	var TemplateJobInstance jbs.Job = &jobs.TemplateJob{}

	TemplateJobInstance.AddDatabaseAccessFunctions(templ.GetDatabaseConnections())
	TemplateJobInstance.AddDatabaseAccessFunctions(templ.GetDatabaseConnections())
	// Lets check the job processing
	err := TemplateJobInstance.Run()
	if err != nil {
		logHandler.ErrorLogger.Println(err)
	}

}
