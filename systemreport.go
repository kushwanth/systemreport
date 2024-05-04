package main

import "fmt"

func main() {
	CPUInfo, err1 := cpuInfo()
	if err1 != nil {
		fmt.Errorf("CPU Info Error")
	}
	fmt.Println(CPUInfo)
}
