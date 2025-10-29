package porting

import (
	"errors"
	"path/filepath"

	"github.com/dubbersthehoser/mayble/internal/csv"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

var importerMap map[string]storage.Importer
var exporterMap map[string]storage.Exporter

func init() {
	importerMap = make(map[string]storage.Importer)
	exporterMap = make(map[string]storage.Exporter)
	importerMap["csv"] = csv.BookLoanCSV{}
	exporterMap["csv"] = csv.BookLoanCSV{}
}

func DriverByFilePath(filePath string) (string, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".csv":
		return "csv", nil
	default:
		return "", errors.New("porting: driver not found for file path.")
	}
}

func ListImporters() []string {
	list := make([]string, len(importerMap))
	count := 0
	for k := range importerMap {
		list[count] = k
		count += 1
	}
	return list
}

func GetImporterByFilePath(filePath string) (storage.Importer, error) {
	driver, err := DriverByFilePath(filePath)
	if err != nil {
		return nil, err
	}
	return GetImporter(driver);
}

func GetImporter(driver string) (storage.Importer, error) {
	impoter, ok := importerMap[driver]
	if !ok {
		return nil, errors.New("import driver not found")
	}
	return impoter, nil
}

func ListExporters() []string {
	list := make([]string, len(exporterMap))
	count := 0
	for k := range exporterMap {
		list[count] = k
		count += 1
	}
	return list
}

func GetExporter(driver string) (storage.Exporter, error) {
	exporter, ok := exporterMap[driver]
	if !ok {
		return nil, errors.New("export driver not found")
	}
	return exporter, nil
}
func GetExporterByFilePath(filePath string) (storage.Exporter, error) {
	driver, err := DriverByFilePath(filePath)
	if err != nil {
		return nil, err
	}
	return GetExporter(driver)
}



