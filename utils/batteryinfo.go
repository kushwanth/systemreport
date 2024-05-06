package utils

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var powerSupplyNameRegex = regexp.MustCompile("/sys/class/power_supply/(.*)/uevent")

// var powerSupplyDeviceNameRegex = regexp.MustCompile("POWER_SUPPLY_NAME=(.*)")
var powerSupplyDeviceTypeRegex = regexp.MustCompile("POWER_SUPPLY_TYPE=(.*)")
var powerSupplyDeviceAConlineRegex = regexp.MustCompile("POWER_SUPPLY_ONLINE=(.*)")
var powerSupplyDeviceBATPresentRegex = regexp.MustCompile("POWER_SUPPLY_PRESENT=(.*)")
var powerSupplyDeviceStatusRegex = regexp.MustCompile("POWER_SUPPLY_STATUS=(.*)")
var powerSupplyDeviceEnergyFullDesignRegex = regexp.MustCompile("POWER_SUPPLY_ENERGY_FULL_DESIGN=(.*)")
var powerSupplyDeviceEnerygyFullNowRegex = regexp.MustCompile("POWER_SUPPLY_ENERGY_FULL=(.*)")
var powerSupplyDeviceLevelRegex = regexp.MustCompile("POWER_SUPPLY_CAPACITY=(.*)")

func getPowerSupplyValue(regexMatches []string) string {
	if len(regexMatches) == 2 {
		if len(regexMatches[1]) > 0 {
			return strings.ToUpper(regexMatches[1])
		}
	}
	return "Unknown"
}

func boolToLabel(boolIntStr string) string {
	label := "Unkown"
	if boolIntStr == "0" {
		label = "No"
	}
	if boolIntStr == "1" {
		label = "Yes"
	}
	return label
}

func GetBatteryInfo() map[string]map[string]string {
	var powerSupplyInfo = map[string]map[string]string{}
	batteryUeventFiles, err1 := filepath.Glob("/sys/class/power_supply/[A-Z]*/uevent")
	if err1 != nil {
		fmt.Errorf("Unable to fetch battery state from /sys/class/power_supply")
		return powerSupplyInfo
	}
	for _, batteryUeventFile := range batteryUeventFiles {
		var batteryInfo = map[string]string{}
		powerSupplyDataBytes, err2 := os.ReadFile(batteryUeventFile)
		if err2 != nil {
			continue
		}
		powerSupplyData := string(powerSupplyDataBytes)
		powerSupplyName := getPowerSupplyValue(powerSupplyNameRegex.FindStringSubmatch(batteryUeventFile))
		//batteryInfo["Label"] = getPowerSupplyValue(powerSupplyDeviceNameRegex.FindStringSubmatch(powerSupplyData))
		batteryInfo["Class"] = getPowerSupplyValue(powerSupplyDeviceTypeRegex.FindStringSubmatch(powerSupplyData))
		if strings.Contains(powerSupplyName, "AC") {
			batteryInfo["Connected"] = boolToLabel(getPowerSupplyValue(powerSupplyDeviceAConlineRegex.FindStringSubmatch(powerSupplyData)))
		}
		if strings.Contains(powerSupplyName, "BAT") {
			batteryInfo["Connected"] = boolToLabel(getPowerSupplyValue(powerSupplyDeviceBATPresentRegex.FindStringSubmatch(powerSupplyData)))
			batteryInfo["State"] = getPowerSupplyValue(powerSupplyDeviceStatusRegex.FindStringSubmatch(powerSupplyData))
			batteryInfo["Power:Charge"] = fmt.Sprintf("%s%%", getPowerSupplyValue(powerSupplyDeviceLevelRegex.FindStringSubmatch(powerSupplyData)))
			energyFullDesign, err3 := strconv.ParseFloat(getPowerSupplyValue(powerSupplyDeviceEnergyFullDesignRegex.FindStringSubmatch(powerSupplyData)), 64)
			energyFullNow, err4 := strconv.ParseFloat(getPowerSupplyValue(powerSupplyDeviceEnerygyFullNowRegex.FindStringSubmatch(powerSupplyData)), 64)
			if err3 != nil || err4 != nil {
				batteryInfo["Power:Capacity"] = "Unknown"
			} else {
				batteryInfo["Power:Capacity"] = fmt.Sprintf("%.0f%%", math.Round((energyFullNow/energyFullDesign)*100))
			}
		}
		powerSupplyInfo[powerSupplyName] = batteryInfo
	}
	return powerSupplyInfo
}
