package main

import (
	"os"
	"runtime/pprof"
	"fmt"
	"os/signal"
	"syscall"
)

func main(){
	f, _ := os.OpenFile("cpu.profile", os.O_CREATE|os.O_RDWR, 0644)
  	defer f.Close()
	pprof.StartCPUProfile(f)
  	defer pprof.StopCPUProfile()
	SetupCloseHandler(f)
	RunSchedule(12)
}
func SetupCloseHandler(f *os.File) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		pprof.StopCPUProfile()
		//f.Close()
		os.Exit(0)
	}()
}
