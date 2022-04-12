package main

import(
	"math"
	"strconv"
)

func GetCFSCpuCores(cpu int)string{
	var tenNum float64 = 0
	for i := 0; i < cpu; i++{
		tenNum += math.Pow(2,float64(i))
	}
	oxNum := strconv.FormatInt(int64(tenNum), 16)
	ox := "0x"
	ox += oxNum
	//fmt.Println(ox)
	return ox
}

func GetFifoCpuSingleCpu(cpu int)string{
	tenNum := math.Pow(2,float64(cpu))
	oxNum := strconv.FormatInt(int64(tenNum), 16)
	ox := "0x"
	ox += oxNum
	return ox
}
