package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:embed data/distro-release-names.json
var distroReleasesData []byte

func getDistroReleaseData() (map[string]string, error) {
	var distroReleases map[string]string
	err := json.Unmarshal(distroReleasesData, &distroReleases)
	if err != nil {
		fmt.Errorf("Unable to read Distro release file")
		return nil, err
	}
	return distroReleases, nil
}

func GetKernelGCCVersion() string {
	procVersionRegex := regexp.MustCompile(`\(([^@]*)\)`)
	procVersionInfo, err1 := os.ReadFile("/proc/version")
	if err1 != nil {
		panic("Unable to read /proc/version")
	}
	gccVersion := "UNKNOWN"
	procVersionMatches := procVersionRegex.FindStringSubmatch(string(procVersionInfo))
	if len(procVersionMatches) == 2 {
		gccVersion = procVersionMatches[1]
	}
	return gccVersion
}

func GetLinuxDistro() string {
	var distroName string
	osReleaseFile, err1 := os.ReadFile("/etc/os-release")
	if err1 == nil {
		osReleaseRegex := regexp.MustCompile("PRETTY_NAME=\"(.*)\"")
		releasePrettyNameMatch := osReleaseRegex.FindStringSubmatch(string(osReleaseFile))
		distroName = strings.TrimSpace(releasePrettyNameMatch[1])
	} else {
		distroReleaseFiles, err1 := filepath.Glob("/etc/*release")
		distroVersionFiles, err2 := filepath.Glob("/etc/*version")
		if err1 != nil || err2 != nil {
			fmt.Errorf("No release Files found")
			distroName = "Linux(Unknown)"
		} else {
			var possibleDistroNames = []string{}
			distroReleases, err3 := getDistroReleaseData()
			if err3 != nil {
				return "UNKNOWN"
			}
			for _, releaseFile := range distroReleaseFiles {
				releaseName, err4 := distroReleases[releaseFile]
				if err4 {
					possibleDistroNames = append(possibleDistroNames, releaseName)
				}
			}
			for _, versionFile := range distroVersionFiles {
				versionName, err5 := distroReleases[versionFile]
				if err5 {
					possibleDistroNames = append(possibleDistroNames, versionName)
				}
			}
			distroName = strings.Join(possibleDistroNames, "")
		}
	}
	return strings.ToUpper(distroName)
}

func GetUptime() string {
	uptimeFile, err1 := os.ReadFile("/proc/uptime")
	if err1 != nil {
		panic("Unable to read /proc/uptime")
	}
	uptimeData := strings.Split(string(uptimeFile), " ")
	uptimeStr, err2 := strconv.ParseFloat(uptimeData[0], 64)
	if err2 != nil {
		fmt.Errorf("Unable to parse uptime")
		return "UNKNOWN"
	}
	var uptime time.Time
	uptime = uptime.Add(time.Duration(uptimeStr) * time.Second)
	parsedUptime := fmt.Sprintf("%02d Hours %02d Minutes %02d Seconds", uptime.Hour(), uptime.Minute(), uptime.Second())
	return parsedUptime
}

func GetKernelInfo() [4]string {
	kernelInfo := [4]string{"UNKNOWN", "UNKNOWN", "UNKNOWN", "UNKNOWN"}
	archFile, err1 := os.ReadFile("/proc/sys/kernel/arch")
	hostnameFile, err2 := os.ReadFile("/proc/sys/kernel/hostname")
	osReleaseFile, err3 := os.ReadFile("/proc/sys/kernel/osrelease")
	kernelVersionFile, err4 := os.ReadFile("/proc/sys/kernel/version")
	if err1 == nil {
		systemArch := strings.TrimSpace(string(archFile))
		kernelInfo[0] = strings.ToLower(systemArch)
	}
	if err2 == nil {
		systemHostname := strings.TrimSpace(string(hostnameFile))
		kernelInfo[1] = strings.ToLower(systemHostname)
	}
	if err3 == nil {
		osRelease := strings.TrimSpace(string(osReleaseFile))
		kernelInfo[2] = strings.ToLower(osRelease)
	}
	if err4 == nil {
		kernelRelease := strings.TrimSpace(string(kernelVersionFile))
		kernelInfo[3] = strings.ToUpper(kernelRelease)
	}
	return kernelInfo
}
