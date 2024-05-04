package main

import (
	"fmt"
	"systemreport/utils"
)

func main() {
	CPUModel, CPUCount := utils.CpuInfo()
	pciDevices := utils.GetAllPCIDevices()
	fmt.Printf("CPU Count %d\nCPU Info %s\n", CPUCount, CPUModel)
	fmt.Printf("PCI Devices\n")
	fmt.Println(pciDevices[1])
	utils.GetLinuxVersion()
}
