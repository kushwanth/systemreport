package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getDistroReleaseData() map[string]string {
	distroReleasePath, _ := filepath.Abs("data/distro-release-names.json")
	file, err := os.ReadFile(distroReleasePath)
	if err != nil {
		fmt.Errorf("Unable to read Distro release file")
	}
	var distroReleases map[string]string
	json.Unmarshal(file, &distroReleases)
	return distroReleases
}

func GetLinuxVersion() (string, string) {
	procVersionRegex := regexp.MustCompile(`Linux version (\S+).*gcc \(GCC\) (\S+)`)
	procVersionInfo, err1 := os.ReadFile("/proc/version")
	if err1 != nil {
		panic("Unable to read /proc/version")
	}
	procVersionMatches := procVersionRegex.FindStringSubmatch(string(procVersionInfo))
	return procVersionMatches[1], procVersionMatches[2]
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
			distroReleases := getDistroReleaseData()
			for _, releaseFile := range distroReleaseFiles {
				releaseName, err3 := distroReleases[releaseFile]
				if err3 {
					possibleDistroNames = append(possibleDistroNames, releaseName)
				}
			}
			for _, versionFile := range distroVersionFiles {
				versionName, err4 := distroReleases[versionFile]
				if err4 {
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
		fmt.Errorf("Unable to read /proc/uptime")
	}
	uptimeData := strings.Split(string(uptimeFile), " ")
	uptimeStr, err2 := strconv.ParseFloat(uptimeData[0], 64)
	if err2 != nil {
		fmt.Errorf("Unablr to parse uptime")
	}
	var uptime time.Time
	uptime = uptime.Add(time.Duration(uptimeStr) * time.Second)
	parsedUptime := fmt.Sprintf("%02d Hours %02d Minutes %02d Seconds", uptime.Hour(), uptime.Minute(), uptime.Second())
	return parsedUptime
}

func GetKernelInfo() [2]string {
	kernelInfo := [2]string{"Unknown", "Unknown"}
	archFile, err1 := os.ReadFile("/proc/sys/kernel/arch")
	hostnameFile, err2 := os.ReadFile("/proc/sys/kernel/hostname")
	if err1 == nil {
		systemArch := strings.TrimSpace(string(archFile))
		kernelInfo[0] = strings.ToUpper(systemArch)
	}
	if err2 == nil {
		systemHostname := strings.TrimSpace(string(hostnameFile))
		kernelInfo[1] = strings.ToUpper(systemHostname)
	}
	return kernelInfo
}
