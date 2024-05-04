package utils

import (
	"os"
	"regexp"
	"slices"
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
var cpuModels []string

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
			modelName := strings.TrimSpace(processorLineData[2])
			cpuProcessor.ModelName = modelName
			if !slices.Contains(cpuModels, modelName) {
				cpuModels = append(cpuModels, modelName)
			}
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

func CpuInfo() (string, int) {
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
	// formattedCPUData, err2 := json.MarshalIndent(cpuProcessors, "", "    ")
	// if err2 != nil {
	// 	return "", 0, err2
	// }
	cpuModel := strings.Join(cpuModels, " ")
	return strings.TrimSpace(cpuModel), len(cpuProcessors)
}
