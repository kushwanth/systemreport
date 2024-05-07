package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PCIID struct {
	VendorID   string `json:"vendorId"`
	VendorName string `json:"vendorName"`
	DeviceID   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`

	SubVendorID   string `json:"subvendorId,omitempty"`
	SubVendorName string `json:"subvendorName,omitempty"`
	SubDeviceID   string `json:"subdeviceId,omitempty"`
	SubDeviceName string `json:"subdeviceName,omitempty"`
}

//go:embed data/pciinfo.json
var pciDevicesData []byte
var pciDevicesIDRegEx = regexp.MustCompile("PCI_ID=(.*)")

func queryPCIInfo(pciIdWithVendorId string) string {
	var pciInfos map[string]string
	json.Unmarshal(pciDevicesData, &pciInfos)
	pciInfo, err := pciInfos[pciIdWithVendorId]
	if !err {
		fmt.Errorf("PCI data doesn't exist")
	}
	return pciInfo
}

func GetAllPCIDevices() map[string]string {
	pciDevices, err1 := filepath.Glob("/sys/bus/pci/devices/*/uevent")
	if err1 != nil {
		panic("Unable to read PCI devices")
	}
	var pciDevicesInfo map[string]string
	pErr := json.Unmarshal(pciDevicesData, &pciDevicesInfo)
	if pErr != nil {
		panic("Unable to fetch PCI Data")
	}
	var pciInfos = map[string]string{}
	for _, pciDevicePath := range pciDevices {
		pciDeviceData, err2 := os.ReadFile(pciDevicePath)
		if err2 != nil {
			fmt.Errorf("Unable to read PCI Ids file")
			continue
		}
		pciIdMatched := pciDevicesIDRegEx.FindStringSubmatch(string(pciDeviceData))
		pciId := strings.TrimSpace(pciIdMatched[1])
		pciDeviceInfo, err3 := pciDevicesInfo[pciId]
		if !err3 {
			fmt.Errorf("PCI device data doesn't exist")
		} else {
			pciInfos[pciId] = pciDeviceInfo
		}
	}
	return pciInfos
}
