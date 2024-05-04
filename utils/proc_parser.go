package utils

import (
	"fmt"
	"os"
	"regexp"
)

func GetLinuxVersion() {
	procVersionRegex := regexp.MustCompile(`Linux version (\S+).*gcc \(GCC\) (\S+)`)
	procVersionInfo, err1 := os.ReadFile("/proc/version")
	if err1 != nil {
		panic("Unable to read /proc/version")
	}
	procVersionMatches := procVersionRegex.FindStringSubmatch(string(procVersionInfo))
	fmt.Println(procVersionMatches[1], procVersionMatches[2])
}
