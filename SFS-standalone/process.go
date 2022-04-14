package main

import (
    //"fmt"

    "github.com/shirou/gopsutil/v3/process"
    // "github.com/shirou/gopsutil/mem"  // to use v2
)

func GetProcessState(pid int) int{
	pid32 := int32(pid)
	v, err := process.NewProcess(pid32)
	if err != nil{
		// terminated
		//fmt.Println("logs error in getting process")
		return 3
	}
	status, _ := v.Status()
	//fmt.Println("logs process'state is ", status)
	if status[0] == "sleep"{
		//fmt.Println("logs sleep")
		return 2
	}else if status[0] == "running"{
		//fmt.Println("logs running")
		return 1
	}
	return 4
}
