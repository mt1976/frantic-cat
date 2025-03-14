package main

import (
	"github.com/mt1976/frantic-cat/app/dao/storage"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/logHandler"
)

func Catalog(cfg *commonConfig.Settings, probeOnly bool) error {
	// This is the main function

	_, err := storage.Catalog(cfg, probeOnly)

	if err != nil {
		logHandler.ErrorLogger.Println("Error exporting storage records: ", err)
		return err
	}
	return nil
}
