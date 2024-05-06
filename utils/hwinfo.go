package utils

import (
	"fmt"
	"os"
	"strings"
)

func getValue(fileBytes []byte, err error) string {
	var value string
	if err != nil {
		fmt.Errorf("Unable to read /sys/devices/virtual/dmi file", err)
		value = "UNKNOWN"
	} else {
		fileData := string(fileBytes)
		value = strings.ReplaceAll(fileData, "\n", "")
	}
	return strings.ToUpper(value)
}

func GetHWInfo() map[string]string {
	var hwInfo = map[string]string{}
	hwInfo["BIOS:Date"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/bios_date"))
	hwInfo["BIOS:Vendor"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/bios_vendor"))
	hwInfo["BIOS:Release"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/bios_release"))
	hwInfo["BIOS:Version"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/bios_version"))
	hwInfo["System:Vendor"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/sys_vendor"))
	hwInfo["Product:Version"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/product_version"))
	hwInfo["Product:SKU"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/product_sku"))
	hwInfo["Product:Name"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/product_name"))
	hwInfo["Product:Family"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/product_family"))
	hwInfo["Board:Name"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/board_name"))
	hwInfo["Board:Vendor"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/board_vendor"))
	hwInfo["Board:Version"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/board_version"))
	hwInfo["Chassis:Vendor"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/chassis_vendor"))
	hwInfo["Chassis:Version"] = getValue(os.ReadFile("/sys/devices/virtual/dmi/id/chassis_version"))
	return hwInfo
}
