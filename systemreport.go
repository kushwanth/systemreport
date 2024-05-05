package main

import (
	"encoding/json"
	"fmt"
	"systemreport/utils"
)

type SystemReport struct {
	OS            string            `json:"OS"`
	KernelVersion string            `json:"Kernel Version"`
	Arch          string            `json:"Architecture"`
	GCCVersion    string            `json:"GCC Version"`
	Hostname      string            `json:"Hostname"`
	CPU           string            `json:"CPU"`
	Threads       int               `json:"Threads"`
	Memory        string            `json:"Memory"`
	Swap          string            `json:"Swap"`
	Uptime        string            `json:"Uptime"`
	Network       []string          `json:"Network"`
	Disk          map[string]string `json:"Disks"`
	Env           map[string]string `json:"OS Env"`
	PCIDevices    map[string]string `json:"PCI Devices"`
}

func main() {
	var systemReport SystemReport
	kernelVersion, gccVersion := utils.GetLinuxVersion()
	distroName := utils.GetLinuxDistro()
	uptime := utils.GetUptime()
	CPUModel, CPUCount := utils.CpuInfo()
	memTotal, memFree, swapTotal, swapFree := utils.MemInfo()
	kernelInfo := utils.GetKernelInfo()
	systemReport.OS = distroName
	systemReport.KernelVersion = kernelVersion
	systemReport.Arch = kernelInfo[0]
	systemReport.GCCVersion = gccVersion
	systemReport.Hostname = kernelInfo[1]
	systemReport.Threads = CPUCount
	systemReport.CPU = CPUModel
	systemReport.Memory = fmt.Sprintf("%d MiB/ %d MiB", memFree, memTotal)
	systemReport.Swap = fmt.Sprintf("%d MiB / %d MiB", swapFree, swapTotal)
	systemReport.Uptime = uptime
	ipAddress := utils.GetIPInfo()
	diskSizes := utils.GetDiskInfo()
	osEnvs := utils.GetOSEnv()
	systemReport.Network = ipAddress
	systemReport.Disk = diskSizes
	systemReport.Env = osEnvs
	pciDevices := utils.GetAllPCIDevices()
	systemReport.PCIDevices = pciDevices
	jsonSystemReport, err := json.MarshalIndent(systemReport, "", "    ")
	if err != nil {
		fmt.Errorf("Unable to generate System Report")
	}
	fmt.Println(string(jsonSystemReport))
}
