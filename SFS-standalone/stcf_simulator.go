package main

import(
	 "os"
        "bufio"
        "strings"
        "strconv"
	"fmt"
	"log"
)
const MAX_RUNNINGTIME = 2147483647

type Exec struct{
	Ac	Action
	Opt	int
}

func Read_optimal(path string) map[int]int{
	//return map[int]int
	file, err := os.Open(path)
        if err != nil{
                log.Fatal(err)
        }

        defer file.Close()
        scanner := bufio.NewScanner(file)
        scanner.Split(bufio.ScanLines)
        var txtlines []string
        for scanner.Scan(){
                txtlines = append(txtlines, scanner.Text())
        }
	dic := make(map[int]int)
	var s []string
	var n int
	var t int
	for _,eachline := range txtlines{
                s = strings.Split(eachline," ")
                n, _ = strconv.Atoi(s[0])
                t, _ = strconv.Atoi(s[1])
		dic[n] = t
        }
	return dic
}

func check_in_select(selected []int, id int) bool{
	for _,v := range selected{
		if v == id{
			return true
		}
	}
	return false
}

func Simulated_execute(workloads []Exec, c_time int, n int)[]Exec {
	//var limit int
	//if c_time >= len(workloads){
	//	limit = len(workloads)
	//}else{
	//	limit = c_time
	//}
	temp := 0
	for k, v := range(workloads){
		if v.Ac.Start > c_time{
			//fmt.Println("logs", "start time", c_time,v.Ac)
			temp = k-1
			break
		}
		temp = len(workloads)
	}
	//if temp == -1{
	//	temp = len(workloads)
	//}
	if temp < 0{
		temp = 0
	}
	running_jobs := workloads[:temp]
	//fmt.Println(c_time, workloads[:temp])
	var selected_id []int
	for i:= 0; i<n; i++{
		max_v := MAX_RUNNINGTIME
		id := -1
		for k,v := range running_jobs {
			if v.Opt < max_v && !check_in_select(selected_id, k){
				max_v = v.Opt
				id = k
			}
		}
		if max_v != MAX_RUNNINGTIME && id != -1{
			//fmt.Println(id, max_v)
			selected_id = append(selected_id, id)
		}
	}
	//fmt.Println(c_time, selected_id)
	//fmt.Println(running_jobs)
	for _, v := range selected_id{
		//fmt.Println(c_time, workloads[v])
		running_jobs[v].Opt -= 1
		if running_jobs[v].Opt <= 0{
			running_jobs[v].Opt = MAX_RUNNINGTIME
			s := strconv.Itoa(c_time - workloads[v].Ac.Start)
			fmt.Println(workloads[v].Ac.JobName, s)
		}
	}
	return workloads
}

func check_finished(workloads []Exec)bool{
	for _,v := range workloads{
		if v.Opt != MAX_RUNNINGTIME{
			return false
		}
	}
	return true
}

func Simulate_schedule(trace []Action, opt_f string, n int){
	var workloads []Exec
	dic := Read_optimal(opt_f)
	for _,v := range trace{
		workloads = append(workloads,Exec{v, dic[v.Para]})
	}
	c_time := 0
	for {
		workloads = Simulated_execute(workloads, c_time, n)
		c_time += 1
		if check_finished(workloads){
			return
		}
	}
}


