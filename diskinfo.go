package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var diskSizeRegex = regexp.MustCompile("/sys/block/(.*)/size")

func GetDiskInfo() (map[string]string, string) {
	var diskSizes = map[string]string{}
	possibleDiskCapacity := 0.0
	diskSizeFiles, err1 := filepath.Glob("/sys/block/*/size")
	if err1 != nil {
		panic("Unable to read /sys/block/*/size")
	}
	for _, diskSizeFile := range diskSizeFiles {
		blockDisk := diskSizeRegex.FindStringSubmatch(diskSizeFile)
		blockDiskName := strings.TrimSpace(blockDisk[1])
		if strings.Contains(blockDiskName, "dm-") || strings.Contains(blockDiskName, "ram") {
			continue
		}
		blockDiskSizeBlocks, err2 := os.ReadFile(diskSizeFile)
		fmtBlockDiskSizeBlocks := strings.Replace(string(blockDiskSizeBlocks), "\n", "", -1)
		blockDiskSize, err3 := strconv.ParseFloat(fmtBlockDiskSizeBlocks, 64)
		possibleDiskCapacity += blockDiskSize
		if err2 != nil || err3 != nil {
			errorOut("Unable to read and parse block size")
		} else {
			diskSize := fmt.Sprintf("%.0f GiB", math.Round(blockDiskSize/(2000000)))
			diskSizes[blockDiskName] = diskSize
		}
	}
	possibleTotalDiskSize := fmt.Sprintf("%.0f GiB", math.Round(possibleDiskCapacity/(2000000)))
	return diskSizes, possibleTotalDiskSize
}
