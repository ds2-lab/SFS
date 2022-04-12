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
)

type Queue struct{
	Core string
	ExecLength int
	LastLayer int
	UpdateValue int
	FirstLayer int
}

var jobs [1000000] int
func handleUDPConnection(conn *net.UDPConn, queue chan int) {
         buffer := make([]byte, 1024)
         n, addr, err := conn.ReadFromUDP(buffer)
         fmt.Println("UDP client : ", addr)
         fmt.Println("Received from UDP client :  ", string(buffer[:n]))
	 now := time.Now()
	 nsec := now.UnixNano()
	 fmt.Println("current time: ",nsec)
         if err != nil {
                 log.Fatal(err)
         }
	 // GET pid and send it to queuei
	 pid_v, _ := strconv.Atoi(string(buffer[:n]))
	 fmt.Println("logs   ", pid_v)
	 if jobs[pid_v] == 0{
		 fmt.Println("Pid ",pid_v, " start")
		 jobs[pid_v] = 1
		 //fmt.Println(reflect.TypeOf(pid_v).String())
		 queue <- pid_v
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

func Listener(queue chan int){
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
		handleUDPConnection(ln,queue)
	}
}

func (q *Queue) Schedule(in chan int, cfs_chan chan int, cpu int){
	on := 1
	s1 := 0
	count := 0
	for {
		select{
		case x := <-in:
			s1 = 0
			fmt.Println("logs path", q.Core, x)
			if on == 0{
				new_pid := -1
				cfs_chan <- new_pid
				on = 1
			}
			UpdateFunc(x, q.Core, "30")
			for s1 < q.ExecLength{
				time.Sleep(time.Duration(1)*time.Millisecond)
				s1 += 1
				if q.CheckTerminated(x) == -1{
					break
				}
			}
			if q.LastLayer != 1{
				go UpdateFunc(x, q.Core, "20")
				cfs_chan <- x
			}else{
				go SwitchFunc(x, GetCFSCpuCores(cpu))
				cfs_chan <- x
			}
		default:
			time.Sleep(time.Duration(1)*time.Millisecond)
			count += 1
			if count >= 2 && on == 1{
				new_pid := -2
				cfs_chan <- new_pid
				on = 0
				count = 0
			}else if on != 1{
				count = 0
			}
		}
	}
}

func HandleCFSChan(in chan int, cfs_value int){
	var a int
	a  = 0
	for{
		select{
		case x:= <-in:
			a += 1
			if x < 0{
				fmt.Println("logs cfs change", cfs_value,x)
			}
		}
	}
}

func (q *Queue) CheckTerminated(job int) int{
	va := 0
	if jobs[job] == 0{
		va = 1
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
	chan1 := make(chan int)
	cfs_chan := make(chan int)
	var queues [1024]Queue
	for i := 0; i < n; i++{
		fmt.Println("logs cpu ", i)
		queues[i] = Queue{GetFifoCpuSingleCpu(i), 20, 1,1,1}
	}
	for i := 0; i < n; i++{
		go queues[i].Schedule(chan1, cfs_chan, n)
	}
	go HandleCFSChan(cfs_chan, 2)
	go Listener(chan1)
}

