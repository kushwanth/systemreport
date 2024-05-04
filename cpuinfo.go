package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Processor struct {
	Id        int64   `json:"Processor"`
	VendorId  string  `json:"Vendor"`
	Model     int64   `json:"Model"`
	ModelName string  `json:"Model Name"`
	MHz       float64 `json:"MHz"`
	CacheSize string  `json:"Cache Size"`
	// Flags     []string `json:"Flags`
	// Bugs      []string `json:"Bugs"`
}

var cpuProcessorDataLineRegEx = regexp.MustCompile("([^:]*?)\\s*:\\s*(.*)$")
var oneLineRegEx = regexp.MustCompile(`\n\s*\n`)

func parseCPUInfo(processorData string) Processor {
	processorDataLines := strings.Split(processorData, "\n")
	cpuProcessor := Processor{}
	for _, processorLine := range processorDataLines {
		processorLineData := cpuProcessorDataLineRegEx.FindStringSubmatch(processorLine)
		switch processorLineData[1] {
		case "processor":
			cpuProcessor.Id, _ = strconv.ParseInt(processorLineData[2], 10, 64)
		case "vendor_id":
			cpuProcessor.VendorId = processorLineData[2]
		case "model":
			cpuProcessor.Model, _ = strconv.ParseInt(processorLineData[2], 10, 64)
		case "model name":
			cpuProcessor.ModelName = processorLineData[2]
		case "cpu MHz":
			cpuProcessor.MHz, _ = strconv.ParseFloat(processorLineData[2], 64)
		case "cache size":
			cpuProcessor.CacheSize = processorLineData[2]
			// case "flags":
			// 	cpuProcessor.Flags = strings.Split(processorLineData[2], " ")
			// case "bugs":
			// 	cpuProcessor.Bugs = strings.Split(processorLineData[2], " ")
		}
	}
	return cpuProcessor
}

func cpuInfo() (string, error) {
	procCPUInfo, err1 := os.ReadFile("/proc/cpuinfo")
	if err1 != nil {
		panic("Unable to read /proc/cpuinfo")
	}
	cpuData := strings.TrimSpace(string(procCPUInfo))
	cpuProcessorsData := oneLineRegEx.Split(cpuData, -1)
	cpuProcessors := []Processor{}
	for _, cpuProcessorData := range cpuProcessorsData {
		cpuProcesserParsedData := parseCPUInfo(cpuProcessorData)
		cpuProcessors = append(cpuProcessors, cpuProcesserParsedData)
	}
	parsedCPUData, err2 := json.MarshalIndent(cpuProcessors, "", "    ")
	if err2 != nil {
		return "", err2
	}
	return string(parsedCPUData), nil
}
