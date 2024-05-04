package utils

import (
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

var pciDevicesIDRegEx = regexp.MustCompile("PCI_ID=[A-Za-z0-9]+:[A-Za-z0-9]+")

func parsePCIIds() []PCIID {
	pciidsPath, _ := filepath.Abs("data/_pciids.json")
	file, err := os.ReadFile(pciidsPath)
	if err != nil {
		panic("Unable to read PCI Ids file")
	}
	var PCIIDs []PCIID
	json.Unmarshal(file, &PCIIDs)
	var PCIInfos = map[string]string{}
	for _, pciId := range PCIIDs {
		pciInfoId := fmt.Sprintf("%s:%s", strings.ToUpper(pciId.VendorID), strings.ToUpper(pciId.DeviceID))
		pciInfoName := fmt.Sprintf("%s %s", pciId.VendorName, pciId.DeviceName)
		PCIInfos[pciInfoId] = pciInfoName
	}
	return PCIIDs
}

func getPCIDataInfo() map[string]string {
	pciInfoPath, _ := filepath.Abs("data/pciinfo.json")
	file, err := os.ReadFile(pciInfoPath)
	if err != nil {
		panic("Unable to read PCI Ids file")
	}
	var pciInfos map[string]string
	json.Unmarshal(file, &pciInfos)
	return pciInfos
}

func queryPCIInfo(pciIdWithVendorId string) string {
	pciInfos := getPCIDataInfo()
	pciInfo, err := pciInfos[pciIdWithVendorId]
	if !err {
		fmt.Errorf("PCI data doesn't exist")
	}
	return pciInfo
}

func GetAllPCIDevices() string {
	pciDevices, err1 := filepath.Glob("/sys/bus/pci/devices/*/uevent")
	if err1 != nil {
		panic("Unable to read PCI devices")
	}
	pciDevicesInfo := getPCIDataInfo()
	var pciInfos = map[string]string{}
	var pciDevicesFormatted []string
	for _, pciDevicePath := range pciDevices {
		pciDeviceData, err2 := os.ReadFile(pciDevicePath)
		if err2 != nil {
			fmt.Errorf("Unable to read PCI Ids file")
			continue
		}
		pciIdMatched := pciDevicesIDRegEx.FindStringSubmatch(string(pciDeviceData))
		pciIds := strings.Split(pciIdMatched[0], "PCI_ID=")
		pciId := strings.TrimSpace(pciIds[1])
		pciDeviceInfo, err3 := pciDevicesInfo[pciId]
		if !err3 {
			fmt.Errorf("PCI device data doesn't exist")
		} else {
			pciInfos[pciId] = pciDeviceInfo
			pciDeviceFormatted := fmt.Sprintf("%s: %s", pciId, pciDeviceInfo)
			pciDevicesFormatted = append(pciDevicesFormatted, pciDeviceFormatted)
		}
	}
	// formattedPCIDevicesData, err4 := json.MarshalIndent(pciInfos, "", "    ")
	// if err4 != nil {
	// 	return "", err4
	// }
	availablePCIDevices := strings.Join(pciDevicesFormatted, "\n")
	return availablePCIDevices
}
