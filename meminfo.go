package main

import (
	"os"
	"strconv"
	"strings"
)

func MemInfo() (int64, int64, int64, int64) {
	memInfoFile, err1 := os.ReadFile("/proc/meminfo")
	if err1 != nil {
		panic("Unable to read /proc/meminfo")
	}
	memInfoStr := strings.TrimSpace(string(memInfoFile))
	memInfoData := strings.Split(memInfoStr, "\n")
	var memInfoKV = map[string]int64{}
	var memTotal, memFree, swapTotal, swapFree int64
	for _, memInfoLine := range memInfoData {
		memData := strings.Split(memInfoLine, ":")
		memDataName := strings.TrimSpace(memData[0])
		memDataValue := strings.Split(memData[1], " kB")
		memKBValue, _ := strconv.ParseInt(strings.TrimSpace(memDataValue[0]), 10, 64)
		memInfoKV[memDataName] = memKBValue / 1024
		switch memDataName {
		case "MemTotal":
			memTotal = memKBValue / 1024
		case "MemFree":
			memFree = memKBValue / 1024
		case "SwapTotal":
			swapTotal = memKBValue / 1024
		case "SwapFree":
			swapFree = memKBValue / 1024
		default:
			continue
		}
	}
	return memTotal, memFree, swapTotal, swapFree
}
