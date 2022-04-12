package main

import(
	"sync"
	//"fmt"
	"time"
	"flag"
	"syscall"
	"fmt"
	//"runtime"
	//"runtime/debug"
)

func main(){
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	//fmt.Println(rLimit)
	rLimit.Max = 1024000
	rLimit.Cur = 1024000
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)	
	}
    	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
    	if err != nil {
        	fmt.Println("Error Getting Rlimit ", err)
    	}
    	fmt.Println("logs Rlimit Final", rLimit)
	var policy string
	flag.StringVar(&policy, "p","m","scheduling policys: m:SFS; c:CFS, s: SRTF")
	var source string
	flag.StringVar(&source, "t","","trace")
	var optimal string
	flag.StringVar(&optimal,"o","optimal.txt","STCF optimal values")
	cpu := flag.Int("n",16,"# of cpu cores")
	fmt.Println("logs main cpu", *cpu)
	flag.Parse()
	fmt.Println("logs main cpu", *cpu)
	flag.Usage()
	if policy == "m"{
		testMLFQ(*cpu, source)
	}else if policy == "c"{
		testCFS(*cpu, source)
	}else{
		testSTCF(*cpu,source, optimal)
	}
}
func testSTCF(cpu int, source string, optimal string){
	//fmt.Println("this is STCF")
	trace, _ := GetTrace(source)
	Simulate_schedule(trace, optimal, cpu)
}
func testMLFQ(cpu int, source string){
	wg := sync.WaitGroup{}
        trace, num := GetTrace(source)
        cache := make(chan PidI)
        wg.Add(1)
        go Scheduler(&wg,cache,cpu,num)
	for i:=0; i < len(trace);i++{
		go Send(trace[i],cache)
		if i < len(trace) - 1{
			time.Sleep(time.Duration(trace[i+1].Start - trace[i].Start)*time.Millisecond)
		}
	}
        //for _, v := range(trace){
        //        go Send(v,cache)
        //}
	//for{
        //        time.Sleep(time.Duration(1000) * time.Millisecond)
	//	debug.FreeOSMemory()
	//	fmt.Println("logs nums of routines", runtime.NumGoroutine())
        //}
        wg.Wait()
}

func testFIFO(cpu int){
	start_time := time.Now()
	wg := sync.WaitGroup{}
        trace, _ := GetTrace("test7")
        cache := make(chan PidI)
        for _, v := range(trace){
                wg.Add(1)
                go ExecuteNoChannel(&wg,v,"F",cache,start_time,"0xff")
        }
        wg.Wait()
}


func testCFS(cpu int, source string){
	start_time := time.Now()
        wg := sync.WaitGroup{}
        trace, _ := GetTrace(source)
        cache := make(chan PidI)
        //go scheduler(&wg,cache)
	cpuC := GetCFSCpuCores(cpu)
	wg.Add(len(trace))
	for i:=0; i < len(trace);i++{
                go ExecuteNoChannel(&wg,trace[i],"N",cache,start_time,cpuC)
		if i < len(trace) - 1{
                	time.Sleep(time.Duration(trace[i+1].Start - trace[i].Start)*time.Millisecond)
		}
        }
	/*
        for _, v := range(trace){
                wg.Add(1)
                go ExecuteNoChannel(&wg,v,"N",cache,start_time,cpuC)
        }
	*/
        wg.Wait()
}

