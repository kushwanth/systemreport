package main

import (
	"net"
	"os"
	"slices"
	"strings"
)

func GetIPInfo() []string {
	var ipAddress []string
	fibTrieFile, err1 := os.ReadFile("/proc/net/fib_trie")
	if err1 != nil {
		errorOut("Unable to read /proc/net/fib_trie")
		return ipAddress
	}
	fibTrieList := strings.Split(string(fibTrieFile), "/32 host LOCAL")
	for _, partFibTire := range fibTrieList {
		partFibTireList := strings.Fields(partFibTire)
		lastElem := partFibTireList[len(partFibTireList)-1]
		ip := net.ParseIP(lastElem)
		if ip != nil && !slices.Contains(ipAddress, ip.String()) {
			ipAddress = append(ipAddress, ip.String())
		}
	}
	ifconfigV6File, err2 := os.ReadFile("/proc/net/if_inet6")
	if err2 != nil {
		errorOut("Unable to read /proc/net/if_inet")
		return ipAddress
	}
	ifconfigV6List := strings.Split(string(ifconfigV6File), "\n")
	for _, ifconfigV6Line := range ifconfigV6List {
		if len(ifconfigV6Line) < 1 {
			continue
		}
		ifconfigV6 := strings.Split(ifconfigV6Line, " ")
		ifconfigV6Str := strings.TrimSpace(ifconfigV6[0])
		ipV6 := []string{}
		for i := 0; i < len(ifconfigV6Str); i += 4 {
			ipV6 = append(ipV6, ifconfigV6Str[i:i+4])
		}
		ipV6Str := strings.Join(ipV6, ":")
		ipAddress = append(ipAddress, ipV6Str)
	}
	return ipAddress
}
