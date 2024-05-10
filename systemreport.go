package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type SystemReport struct {
	OS            string                       `json:"OS"`
	KernelRelease string                       `json:"Kernel:Release"`
	KernelVersion string                       `json:"Kernel:Version"`
	Arch          string                       `json:"Architecture"`
	GCCVersion    string                       `json:"GCC:Version"`
	Hostname      string                       `json:"Hostname"`
	CPU           string                       `json:"CPU"`
	GPU           []string                     `json:"GPU"`
	Threads       string                       `json:"Threads"`
	Memory        string                       `json:"Memory"`
	Swap          string                       `json:"Swap"`
	Uptime        string                       `json:"Uptime"`
	Network       []string                     `json:"Network"`
	Disk          map[string]string            `json:"Disks"`
	Env           map[string]string            `json:"OS:Env"`
	Power         map[string]map[string]string `json:"Power"`
	System        map[string]string            `json:"System"`
	PCIDevices    map[string]string            `json:"PCI:Devices"`
}

func getSystemReport() SystemReport {
	var systemReport SystemReport
	gccVersion := GetKernelGCCVersion()
	distroName := GetLinuxDistro()
	uptime := GetUptime()
	CPUModel, CPUCount := CpuInfo()
	memTotal, memFree, swapTotal, swapFree := MemInfo()
	kernelInfo := GetKernelInfo()
	systemReport.OS = distroName
	systemReport.KernelRelease = kernelInfo[2]
	systemReport.KernelVersion = kernelInfo[3]
	systemReport.Arch = kernelInfo[0]
	systemReport.GCCVersion = gccVersion
	systemReport.Hostname = kernelInfo[1]
	systemReport.Threads = strconv.Itoa(CPUCount)
	systemReport.CPU = CPUModel
	systemReport.Memory = fmt.Sprintf("%d MB / %d MB", memFree, memTotal)
	systemReport.Swap = fmt.Sprintf("%d MB / %d MB", swapFree, swapTotal)
	systemReport.Uptime = uptime
	ipAddress := GetIPInfo()
	diskSizes := GetDiskInfo()
	osEnvs := GetOSEnv()
	systemReport.Network = ipAddress
	systemReport.Disk = diskSizes
	systemReport.Env = osEnvs
	power := GetBatteryInfo()
	system := GetHWInfo()
	systemReport.Power = power
	systemReport.System = system
	pciDevices, gpuInfo := GetAllPCIDevices()
	if len(gpuInfo) == 0 {
		systemReport.GPU = []string{"UNKNOWN"}
	} else {
		systemReport.GPU = gpuInfo
	}
	systemReport.PCIDevices = pciDevices
	return systemReport
}

func main() {
	systemReport := getSystemReport()
	stringOut("OS", systemReport.OS)
	stringOut("Kernel:Release", systemReport.KernelRelease)
	stringOut("Architecture", systemReport.Arch)
	stringOut("GCC:Version", systemReport.GCCVersion)
	stringOut("CPU", systemReport.CPU)
	arrayOut("GPU", systemReport.GPU)
	stringOut("Threads", systemReport.Threads)
	stringOut("Memory", systemReport.Memory)
	arrayOut("Network", systemReport.Network)
	stringOut("Uptime", systemReport.Uptime)
	mapOut("LANG", systemReport.Env)
	mapOut("SHELL", systemReport.Env)
	mapOut("USER", systemReport.Env)
	mapOut("System:Vendor", systemReport.System)
	mapOut("Product:Family", systemReport.System)
	mapOut("BIOS:Version", systemReport.System)
	mapOut("XDG_BACKEND", systemReport.Env)
	stringOut("PCI:Devices", strconv.Itoa(len(systemReport.PCIDevices)))
	jsonSystemReportBytes, err := json.MarshalIndent(systemReport, "", " ")
	if err != nil {
		errorOut("Unable to generate System Report")
	} else {
		os.WriteFile("/tmp/systemreport.json", jsonSystemReportBytes, 0644)
		warningOut("Full System Report stored at /tmp/systemreport.json")
	}
}
