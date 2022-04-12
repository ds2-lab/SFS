package main

import(
    "fmt"
    "strconv"
    //"log"
    "os/exec"
    "time"
    "sync"
    "log"
    "syscall"
    //"sort"
)

var CFS_int int64 = 4
var jobs [1000000] int64
var credits [1000000] int
var remain [100000] float64
var pids [10000000] int
var et [1000000] time.Time

func receive(in chan PidI, queue chan PidI, core string,wg *sync.WaitGroup, num int, ts_chan chan PidI, ts *Threshold){
	defer wg.Done()
	//receiver 1)  
	//         2) delete jobs if receive the job again
	//         3) send job to first queue
	//fmt.Println("logs receive cpu", core)
	num_job := 0
	var init_credit int
	for{
		select{
		case x := <-in:
			if jobs[x.Id] == 0 && x.Credit == -3{
				jobs[x.Id] = 1
				if ts.T > 6{
					init_credit = ts.T
				}else{
					init_credit = 6
				}
				new_x := PidI{x.Pid, x.Job, x.N, x.Id,x.St, init_credit}
				credits[x.Id] = init_credit
				go Execute(new_x, "F", in, core, queue)
				ts_chan <- new_x
			}else if jobs[x.Id] == 3 && credits[x.Id] > 0{
				//sleep & wake jobs
                                fmt.Println("logs this is sleep & waitup jobs")
                                jobs[x.Id] = 2
                                cur_credit := int(remain[x.Id]*float64(ts.T))
                                new_x := PidI{x.Pid, x.Job, x.N, x.Id,time.Now(), cur_credit}
                                ts_chan <- new_x
				queue <- new_x
			}else{
				jobs[x.Id] = 0
				credits[x.Id] = -2
				num_job += 1
				if (num_job >= num){
					return
				}
			}
			/**
			else{
				//sleep & wake jobs
				fmt.Println("logs this is sleep & waitup jobs")
				jobs[x.Id] = 3
				cur_credit := credits[x.Id]
				new_x := PidI{x.Pid, x.Job, x.N, x.Id,x.St, cur_credit}
				ts_chan <- new_x
				//queue <- new_x
			}
			**/
		}
	}
        //default:
                //continue
}

type Queue struct{
	Core string
	ExecLength	int
	LastLayer	int
	UpdateValue	int
	FirstLayer	int
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

// boost sleep jobs

func boostSleepingJobs(in chan PidI){
	for{
		for k, v := range jobs{
			if v == 2{
				fmt.Println("logs jobs k, v ", k, v, pids[k])
				if GetProcessState(pids[k]) == 1{
					//PidI{0,job.Job,job.N,job.Id,time.Now(), job.Credit}
					fmt.Println("logs sleeping jobs activate and send")
					new_pid := PidI{pids[k], "fib", 20, k, time.Now(),credits[k]}
					jobs[k] = 3
					in <- new_pid
				}
			}
		}
		time.Sleep(time.Duration(1)*time.Millisecond)
	}
}
// cfs boost policy
func boostCFSJobs(in chan PidI, threshold int, ts_chan chan PidI){
        for{
                for k, v := range jobs{
                        if v == 4{
				o := time.Now()
                                //fmt.Println("logs jobs k, v ", k, v)
				if int(o.Sub(et[k]).Milliseconds()) > threshold && GetProcessState(pids[k]) == 1{
                                        //PidI{0,job.Job,job.N,job.Id,time.Now(), job.Credit}
                                        fmt.Println("logs cfs job boost")
                                        new_pid := PidI{pids[k], "fib", 20, k, time.Now(),credits[k]}
                                        jobs[k] = 5
                                        in <- new_pid
					ts_chan <- new_pid
                                }else if GetProcessState(pids[k]) ==3{
					jobs[k] = 6
				}
                        }
                }
                time.Sleep(time.Duration(1)*time.Millisecond)
        }
}
// threshold policy
func (t *Threshold) AdjustThreshold(ts_chan chan PidI, period int, n int){
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
				interval_array[count-1] = interval_time
				count = 0
				fmt.Println("interval_time", interval_array)
				t.T = calcuMean(interval_array)*n
				fmt.Println("logs New threshold ", t.T)
				interval_array = make([]int, period)
			}else{
				inc_time := time.Now()
				interval_time = int(inc_time.Sub(cur_time).Milliseconds())
				fmt.Println("logs iat ",interval_time)
				interval_array[count-1] = interval_time
				cur_time = inc_time
			}

		}
	}
}
//var CFS_int int64 = 2

type RWMap struct {
    sync.RWMutex
    m map[string]PidI
}

func (q *Queue) CheckTerminated(job PidI, actions RWMap)int{
	va := 0
	if jobs[job.Id] == 0{
		va = -1
	}else{
		va = 1
	}
	return va
}

func SwitchFunc(pid int, core string){
	var cmd *exec.Cmd
	cmd = exec.Command("schedtool","-N", "-a", core, strconv.Itoa(pid))
	err := cmd.Start()
	if err != nil{
		log.Fatal(err)
	}
	cmd.Wait()

}

func UpdateFunc(pid int, core string, p string){
	var cmd *exec.Cmd
	cmd = exec.Command("schedtool","-F", "-p", p, "-a", core, strconv.Itoa(pid))
	err := cmd.Start()
        if err != nil{
                log.Fatal(err)
        }
	cmd.Wait()
}


func UpdateCFScore(direct int, cfs_value int64, update_v int)int64{
	if direct == -1{
		return cfs_value - int64(update_v)
	}else{
		return cfs_value + int64(update_v)
	}
}


func (q *Queue) Schedule(actions RWMap, cache chan PidI, in chan PidI, out chan PidI, cfs_chan chan PidI, cpu int, ts *Threshold){
	on := 1
	count := 0
	s1 := 0
	for {
		select{
		//receive jobs from prev layer
		case x, _ := <-in:
			pids[x.Id] = x.Pid
			if q.FirstLayer == 1{
				fmt.Println("logs q1 Time start", x.Job, time.Now())
			}
			s1 = 0
			fmt.Println("logs path", q.Core, x)
			if on == 0{
				new_pid := PidI{-1, "minus", q.UpdateValue, -1,time.Now(), x.Credit}
                                cfs_chan <- new_pid
				on = 1
			}
			o := time.Now()
                        queue_delay := int(o.Sub(x.St).Milliseconds())
                        fmt.Println("logs queue delay ", x.Id, queue_delay)
                        // use default cfs scheduleor
                        if int(o.Sub(x.St).Milliseconds()) > 3 * ts.T{
                                jobs[x.Id] = 2
                                credits[x.Id] = x.Credit
                                cfs_chan <-x
                                pids[x.Id] = x.Pid
                                et[x.Id] = time.Now()
                                fmt.Println("logs cfs scheduler running")
                                continue
                        }
			//actions.Lock()
			UpdateFunc(x.Pid, q.Core, "30")
			//actions.Unlock()
			exec_time := 0
			if ts.T == 0{
				exec_time = 6
			}else{
				exec_time = ts.T
			}
			if credits[x.Id] > 0{
				exec_time = credits[x.Id]
			}
			for s1 < exec_time{
				//fmt.Println("logs exec", q.Core, s1)
				time.Sleep(time.Duration(1)*time.Millisecond)
				s1 += 1
				if(q.CheckTerminated(x,actions) == -1){
					break
				}else if GetProcessState(x.Pid) == 2{
					jobs[x.Id] = 2
					credits[x.Id] =x.Credit - s1
					remain[x.Id] = float64(x.Credit - s1)
					break
				}
				if s1 >= x.Credit{
					break
				}
			}
			//if(q.CheckTerminated(x,actions) == 1){
				//transfer job to next later
			if (q.LastLayer != 1){
				go UpdateFunc(x.Pid, q.Core, "20")
				//SwitchFunc(x.Pid, GetCFSCpuCores(cpu))
				out <- x
			}else{
                               	//UpdateFunc(x.Pid, q.Core, "20")
				go SwitchFunc(x.Pid, GetCFSCpuCores(cpu))
				cfs_chan <- x
				if jobs[x.Id] != 2{
					//4 means cfs state
					jobs[x.Id] = 4
				}
				pids[x.Id] = x.Pid
				et[x.Id] = time.Now()
			}
			//}
		default:
			time.Sleep(time.Duration(1) * time.Millisecond)
			count += 1
			if count >= 2 && on == 1{
				new_pid := PidI{-1, "plus", q.UpdateValue,-1,time.Now(), 20}
				cfs_chan <- new_pid
				on = 0
				count = 0
			}else if on != 1{
				count = 0
			}
		}
	}
}

func HandleCFSChan(actions RWMap, in chan PidI, m map[string]int, cfs_value int64){
	//do nothing
	var a int64 = 0
	for{
		select{
		case x,_ := <-in:
			a += 1
			if a >= 1000000 && x.Pid == -1{
				fmt.Println(a)
			}
		}
	}
}

func Scheduler(wg *sync.WaitGroup, cache chan PidI, cpu int, num int){
	defer wg.Done()
	wg_receive := sync.WaitGroup{}
	var rLimit syscall.Rlimit
        err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
        if err != nil {
                fmt.Println("Error Getting Rlimit ", err)
        }
	tsChan := make(chan PidI)
	chan1 := make(chan PidI)
        chan2 := make(chan PidI)
        //chan3 := make(chan PidI)
        cfs_chan := make(chan PidI)
        con_map := make(map[string]int)
        actions := make(map[string]PidI)
        con_actions := RWMap{m:actions}
	ts_instance := Threshold{20}
        //layer 1
	var queues [1024]Queue
	for i := 0; i < cpu; i++{
		fmt.Println("logs cpu", i)
		queues[i] = Queue{GetFifoCpuSingleCpu(i),20,1,1,1}
	}
	for i:= 0; i < cpu; i++{
		fmt.Println("logs cpu", i)
		go queues[i].Schedule(con_actions, cache, chan1, chan2, cfs_chan,cpu,&ts_instance)
	}
	go HandleCFSChan(con_actions, cfs_chan, con_map, int64(2))
	wg_receive.Add(1)
	go receive(cache, chan1, GetCFSCpuCores(cpu), &wg_receive, num, tsChan, &ts_instance)
	go ts_instance.AdjustThreshold(tsChan, 200, cpu)
	go boostSleepingJobs(chan1)
	boostCFSJobs(chan1, 20000, tsChan)
	wg_receive.Wait()

}


