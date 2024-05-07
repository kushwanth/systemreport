package main

import (
	"os"
	"regexp"
	"slices"
)

var osEnvLineRegEx = regexp.MustCompile("^(.*)=(.*)$")
var necessaryOSEnvs = []string{"USER", "SHELL", "XDG_BACKEND", "XDG_SESSION_DESKTOP", "XDG_SESSION_TYPE", "LANG", "XDG_CURRENT_DESKTOP"}

func GetOSEnv() map[string]string {
	var osEnv = map[string]string{}
	osEnvData := os.Environ()
	for _, osEnvLine := range osEnvData {
		osEnvMatch := osEnvLineRegEx.FindStringSubmatch(osEnvLine)
		if slices.Contains(necessaryOSEnvs, osEnvMatch[1]) {
			osEnv[osEnvMatch[1]] = osEnvMatch[2]
		}
	}
	return osEnv
}
