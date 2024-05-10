package main

import (
	"os"
	"regexp"
	"slices"
	"strings"
)

var osEnvLineRegEx = regexp.MustCompile("^(.*)=(.*)$")
var necessaryOSEnvs = []string{"USER", "XDG_BACKEND", "XDG_SESSION_DESKTOP", "XDG_SESSION_TYPE", "LANG", "XDG_CURRENT_DESKTOP"}

func getShellName(shell string) string {
	shellName := shell
	if strings.Contains(shell, "bash") {
		shellName = "BASH"
	} else if strings.Contains(shell, "sh") {
		shellName = "SH"
	} else if strings.Contains(shell, "zsh") {
		shellName = "ZSH"
	} else if strings.Contains(shell, "fish") {
		shellName = "FISH"
	} else if strings.Contains(shell, "csh") {
		shellName = "CSH"
	} else if strings.Contains(shell, "tcsh") {
		shellName = "TCSH"
	} else {
		shellName = shell
	}
	return shellName
}

func GetOSEnv() map[string]string {
	var osEnv = map[string]string{}
	osEnvData := os.Environ()
	for _, osEnvLine := range osEnvData {
		osEnvMatch := osEnvLineRegEx.FindStringSubmatch(osEnvLine)
		if osEnvMatch[1] == "SHELL" {
			osEnv["SHELL"] = getShellName(strings.ToLower(osEnvMatch[2]))
		} else if slices.Contains(necessaryOSEnvs, osEnvMatch[1]) {
			osEnv[osEnvMatch[1]] = strings.ToLower(osEnvMatch[2])
		} else {
			continue
		}
	}
	return osEnv
}
