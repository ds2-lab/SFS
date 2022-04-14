package main

import(
	"time"
	"fmt"
	"log"
	"sync"
	"strconv"
	"os/exec"
)

type PidI struct{
	Pid	int
	Job	string
	N	int
	Id	int
	St	time.Time
	Credit	int
}

func Send(job Action, pids chan PidI){
	// Send just send request to receiver
	o := time.Now()
	//time.Sleep(time.Duration(job.Start)*time.Millisecond)
	new_pid := PidI{-10, job.JobName, job.Para,job.Id,o, -3}
	pids <- new_pid
}
func Execute(job PidI, p string, pids chan PidI, core string, queue chan PidI){
	// execute request and also update job direction
	var cmd *exec.Cmd
	start_time := job.St
	t1 := time.Now()
	if p == "N"{
		cmd = exec.Command("schedtool","-N","-a",core,"-e","python", "fib.py", strconv.Itoa(job.N),strconv.Itoa(job.Id))
	}else{
		//cmd = exec.Command("schedtool","-N","-a",core,"-e","python","fib.py", strconv.Itoa(job.N))
		cmd = exec.Command("schedtool","-F","-p","20","-a",core,"-e","python","fib.py", strconv.Itoa(job.N),strconv.Itoa(job.Id))
	}
	err := cmd.Start()
	if err != nil{
		log.Fatal("logs exec 1", err)
	}
	tw := time.Now()
	fmt.Println("logs wait time",tw.Sub(t1))
	//actions.m[job.Job] = cmd.Process.Pid
	//new_pid := PidI{0,job.Job,job.N,job.Id}
	pid := cmd.Process.Pid
	var new_pid PidI
	if cmd != nil{
		new_pid = PidI{pid,job.Job,job.N,job.Id,time.Now(), job.Credit}
	}else{
		new_pid = PidI{0,job.Job,job.N,job.Id,time.Now(), job.Credit}
	}
	//actions.Lock()
	//actions.m[job.Job] = new_pid
	//actions.Unlock()
	queue <- new_pid
	err = cmd.Wait()
	if err != nil{
		log.Fatal("exec 2",err)
	}
	t2 := time.Now()
	new_pid.Credit = -2
	pids <- new_pid
	fmt.Println(job.Job,t2.Sub(t1).Milliseconds())
	//cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
	//err = cmd.Start()
	//if err != nil{
	//	log.Fatal("logs exec 3", err)
	//}
	//debug.FreeOSMemory()
        //log4sys.Warn("NumGoroutine:",runtime.NumGoroutine())
	fmt.Println("logs TIME: ",job.Job, t1.Sub(start_time), t2.Sub(start_time))
}

func ExecuteNoChannel(wg *sync.WaitGroup, job Action, p string, pids chan PidI, start_time time.Time, cpuC string){
        defer wg.Done()
        //time.Sleep(time.Duration(job.Start) * time.Millisecond)
        t1 := time.Now()
        var cmd *exec.Cmd
        if p == "N"{
		cmd = exec.Command("schedtool","-N","-a",cpuC,"-e","python", job.Exec, strconv.Itoa(job.Para),strconv.Itoa(job.Id))
        }else{
                cmd = exec.Command("schedtool","-R","-p","20","-a","0x1","-e","python", job.Exec, strconv.Itoa(job.Para),strconv.Itoa(job.Id))
        }
	err := cmd.Start()
	if err != nil{
		log.Fatal("exec 1", err)
	}
	tw := time.Now()
        fmt.Println("logs wait time",tw.Sub(t1))
	//pid := cmd.Process.Pid
        //new_pid := PidI{cmd.Process.Pid,job.JobName}
        //pids <- new_pid
        //err := cmd.Wait()
        //if err != nil{
        //        log.Fatal(err)
        //}
	err = cmd.Wait()
        if err != nil{
                log.Fatal("exec 2",err)
        }

        t2 := time.Now()
        //pids <- new_pid
        fmt.Println(job.JobName,t2.Sub(t1).Milliseconds())
	fmt.Println("logs TIME: ",job.JobName, t1.Sub(start_time), t2.Sub(start_time))
}

