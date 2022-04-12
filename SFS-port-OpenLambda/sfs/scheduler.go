package main

import (
	"fmt"
	"net"
	"log"
	"time"
	//"syscall"
	"os/exec"
	"strconv"
	//"reflect"
	"encoding/json"
)

type Queue struct{
	Core string
	ExecLength int
	LastLayer int
	UpdateValue int
	FirstLayer int
}

var jobs [1000000] int
var pids [10000000] int
var et [1000000] time.Time
var credits [1000000] int
type Job struct{
	id	int
	credit  int
	pid	int
	st	time.Time
}

func handleUDPConnection(conn *net.UDPConn, queue chan Job, ts_chan chan Job, ts *Threshold) {
         buffer := make([]byte, 1024)
         n, addr, err := conn.ReadFromUDP(buffer)
         fmt.Println("UDP client : ", addr)
         fmt.Println("Received from UDP client :  ", string(buffer[:n]))
         if err != nil {
                 log.Fatal(err)
         }
	 // GET pid and send it to queue
	 var result map[string]string
	 json.Unmarshal(buffer[:n], &result)
	 pid_v,_ := strconv.Atoi(result["pid"])
	 id,_ := strconv.Atoi(result["id"])
	 fmt.Println("pidv: ", pid_v)
	 fmt.Println("logs   ", pid_v)
	 var init_credit int
	 if jobs[id] == 0{
		 fmt.Println("Pid ",pid_v, " start")
		 jobs[id] = 1
		 pids[id] = pid_v
		 //fmt.Println(reflect.TypeOf(pid_v).String())
		 if ts.T > 6{
			init_credit = ts.T
		 }else{
			init_credit = 6
		 }
		 new_job := Job{id, init_credit, pid_v, time.Now()}
		 credits[id] = init_credit
		 queue <- new_job
		 ts_chan <- new_job

		 fmt.Println("queue send")
	 }else{
		jobs[pid_v] = 0
		fmt.Println("Pid ",pid_v, " end")
	}
	message := []byte("Hello UDP client!")
        _, err = conn.WriteToUDP(message, addr)
	if err != nil {
                log.Println(err)
        }
}

func Listener(queue chan Job, ts_chan chan Job, ts *Threshold){
	hostName := "localhost"
	portNum := "4009"
	service := hostName + ":" + portNum
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	if err != nil{
		 log.Fatal(err)
	}
	//setup listener for incoming udp connection
	ln, err := net.ListenUDP("udp",udpAddr)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("UDP server up and listening on port 4009")
	defer ln.Close()

	for{
		handleUDPConnection(ln,queue, ts_chan,ts)
	}
}

type Threshold struct{
	T	int
}

func calcuMean(n []int)int{
	total := 0
	for _, v := range n{
		total += v
	}
	return total/len(n)
}
func boostSleepingJobs(in chan Job, ts_chan chan Job){
	for{
		for k, v := range jobs{
			if v == 2{
				fmt.Println("logs jobs k, v ", k, v, pids[k])
				if GetProcessState(pids[k]) == 1{
					//PidI{0,job.Job,job.N,job.Id,time.Now(), job.Credit}
					fmt.Println("logs sleeping jobs activate and send")
					new_pid := Job{k, credits[k], pids[k],time.Now()}
					jobs[k] = 3
					in <- new_pid
					ts_chan <- new_pid

				}
			}
		}
		time.Sleep(time.Duration(8)*time.Millisecond)
	}
	//time.Sleep(time.Duration(1)*time.Millisecond)
}

func (t *Threshold) AdjustThreshold(ts_chan chan Job, period int, n int){
	//st := time.Now()
	cur_time := time.Now()
	var interval_time int
	count := 0
	interval_array := make([]int, period)
	for {
		select{
			case _ = <-ts_chan:
			count += 1
			if count >= period{
				inc_time := time.Now()
                                interval_time = int(inc_time.Sub(cur_time).Milliseconds())
				//interval_array = append(interval_array, interval_time)
				interval_array[count-1] = interval_time
				//fmt.Println("logs interval array 1", interval_array)
				count = 0
				//et := time.Now()
				//period_time := et.Sub(st).Milliseconds()
				//st = et
				//sort.Ints(interval_array)
				//t.T = calcuMean(interval_array)*n
				fmt.Println("interval_time", interval_array)
				//t.T = interval_array[100]*n
				t.T = calcuMean(interval_array)*n
				//t.T = (n * int(period_time)) / period
				fmt.Println("logs New threshold ", t.T)
				interval_array = make([]int, period)
				//fmt.Println("logs interval array 2", interval_array)
			}else{
				inc_time := time.Now()
				interval_time = int(inc_time.Sub(cur_time).Milliseconds())
				fmt.Println("logs iat ",interval_time)
				//interval_array = append(interval_array, interval_time)
				interval_array[count-1] = interval_time
				cur_time = inc_time
			}

		}
	}
}

func (q *Queue) Schedule(in chan Job, cfs_chan chan Job, cpu int,  ts *Threshold){
	s1 := 0
	for {
		select{
		case x := <-in:
			pids[x.id] = x.pid
			if jobs[x.id] == 3{
				jobs[x.id] = 1
			}
			o := time.Now()
			queue_delay := int(o.Sub(x.st).Milliseconds())
			if queue_delay > 3*ts.T{
				jobs[x.id] = 4
				credits[x.id] = x.credit
				cfs_chan <- x
				et[x.id] = time.Now()
				continue
			}
			UpdateFunc(x.pid, q.Core, "30")
			exec_time := 0
			if ts.T == 0{
				exec_time = 6
			}else{
				exec_time = ts.T
			}
			if credits[x.id] > 0{
				exec_time = credits[x.id]
			}
			for s1 < exec_time{
				time.Sleep(time.Duration(1)*time.Millisecond)
				s1 += 1
				if q.CheckTerminated(x.pid) == -1{
					break
				}else if GetProcessState(x.pid) == 2{
					jobs[x.id] = 2
					credits[x.id] = x.credit - s1
					break
				}
				if s1 > x.credit{
					break
				}
			}
			if q.LastLayer != 1{
				go UpdateFunc(x.pid, q.Core, "20")
				cfs_chan <- x
			}else{
				go SwitchFunc(x.pid, GetCFSCpuCores(cpu))
				cfs_chan <- x
				if jobs[x.id] != 2{
					//4 means cfs state
					jobs[x.id] = 4
				}
				pids[x.id] = x.pid
				et[x.id] = time.Now()
			}
		default:
			time.Sleep(time.Duration(1)*time.Millisecond)
		}
	}
}

func HandleCFSChan(in chan Job, cfs_value int){
	var a int
	a  = 0
	for{
		select{
		case x:= <-in:
			a += 1
			if x.pid < 0{
				fmt.Println("logs cfs change", cfs_value,x)
			}
		}
	}
}

func (q *Queue) CheckTerminated(job int) int{
	va := 0
	if jobs[job] == 0{
		va = -1
	}else{
		va = 1
	}
	return va
}

func SwitchFunc(pid int, core string){
	var cmd *exec.Cmd
	cmd = exec.Command("schedtool","-N","-a",core, strconv.Itoa(pid))
	err := cmd.Start()
	if err != nil{
		log.Fatal(err)
	}
	cmd.Wait()
}

func UpdateFunc(pid int, core string, p string){
	var cmd *exec.Cmd
	cmd = exec.Command("schedtool","-F","-p",p,"-a",core, strconv.Itoa(pid))
	err := cmd.Start()
	if err != nil{
		log.Fatal(err)
	}
	cmd.Wait()
}

func RunSchedule(n int){
	fmt.Println(n)
	chan1 := make(chan Job)
	cfs_chan := make(chan Job)
	tsChan := make(chan Job)
	var queues [1024]Queue
	 ts_instance := Threshold{20}
        go ts_instance.AdjustThreshold(tsChan, 100, n)
	for i := 0; i < n; i++{
		fmt.Println("logs cpu ", i)
		queues[i] = Queue{GetFifoCpuSingleCpu(i), 20, 1,1,1}
	}
	for i := 0; i < n; i++{
		go queues[i].Schedule(chan1, cfs_chan, n,&ts_instance)
	}
	go HandleCFSChan(cfs_chan, 2)
	go Listener(chan1,tsChan,&ts_instance)
	go boostSleepingJobs(chan1,tsChan)
	for{
		time.Sleep(time.Duration(1)*time.Second)
	}
}

