package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PCIID struct {
	VendorID      string `json:"vendorId"`
	VendorName    string `json:"vendorName"`
	DeviceID      string `json:"deviceId"`
	DeviceName    string `json:"deviceName"`
	SubVendorID   string `json:"subvendorId,omitempty"`
	SubVendorName string `json:"subvendorName,omitempty"`
	SubDeviceID   string `json:"subdeviceId,omitempty"`
	SubDeviceName string `json:"subdeviceName,omitempty"`
}

//go:embed data/pciinfo.json
var pciDevicesData []byte
var pciDevicesIDRegEx = regexp.MustCompile("PCI_ID=(.*)")
var pciModaliasRegEx = regexp.MustCompile("pci:v0000([A-Za-z0-9]*)d0000([A-Za-z0-9]*)sv0000([A-Za-z0-9]*)sd0000([A-Za-z0-9]*)bc([A-Za-z0-9]*)sc([A-Za-z0-9]*)i([A-Za-z0-9]*)")

func queryPCIInfo(pciIdWithVendorId string) string {
	var pciInfos map[string]string
	json.Unmarshal(pciDevicesData, &pciInfos)
	pciInfo, err := pciInfos[pciIdWithVendorId]
	if !err {
		errorOut("PCI data doesn't exist")
	}
	return pciInfo
}

func GetAllPCIDevices() (map[string]string, []string) {
	pciModaliasFiles, err1 := filepath.Glob("/sys/bus/pci/devices/*/modalias")
	if err1 != nil {
		panic("Unable to read PCI devices")
	}
	var pciDevicesInfo map[string]string
	pErr := json.Unmarshal(pciDevicesData, &pciDevicesInfo)
	if pErr != nil {
		panic("Unable to fetch PCI Data")
	}
	var pciInfos = map[string]string{}
	var gpuInfo []string
	for _, pciModaliasFile := range pciModaliasFiles {
		pciModaliasData, err2 := os.ReadFile(pciModaliasFile)
		if err2 != nil {
			errorOut("Unable to read PCI modalias file")
			continue
		}
		pciModaliasMatched := pciModaliasRegEx.FindStringSubmatch(string(pciModaliasData))
		pciId := strings.Join(pciModaliasMatched[1:3], ":")
		pciDeviceInfo, err3 := pciDevicesInfo[pciId]
		if err3 {
			if strings.Contains(pciModaliasMatched[5], "03") {
				gpuInfo = append(gpuInfo, pciDeviceInfo)
			}
			pciCustomId := pciModaliasMatched[5] + ":" + pciId
			pciInfos[pciCustomId] = pciDeviceInfo
		}
	}
	return pciInfos, gpuInfo
}
