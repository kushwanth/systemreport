package main

import (
	"encoding/json"
	"fmt"
)

type SystemReport struct {
	OS            string                       `json:"OS"`
	KernelRelease string                       `json:"Kernel:Release"`
	KernelVersion string                       `json:"Kernel:Version"`
	Arch          string                       `json:"Architecture"`
	GCCVersion    string                       `json:"GCC:Version"`
	Hostname      string                       `json:"Hostname"`
	CPU           string                       `json:"CPU"`
	Threads       int                          `json:"Threads"`
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

func main() {
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
	systemReport.Threads = CPUCount
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
	pciDevices := GetAllPCIDevices()
	systemReport.PCIDevices = pciDevices
	jsonSystemReport, err := json.MarshalIndent(systemReport, "", "    ")
	if err != nil {
		fmt.Errorf("Unable to generate System Report")
	}
	fmt.Println(string(jsonSystemReport))

}
