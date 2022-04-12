package main

import(
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
)

type Action struct{
	JobName	string
	Exec	string
	Para	int
	Start	int
	Id	int
}


func GetTrace(path string)([]Action, int){
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
	trace := []Action{}
	var s []string
	var i int
	var f int
	var id int
	var newAction Action
	var num int = 0
	for _,eachline := range txtlines{
		s = strings.Split(eachline," ")
		i, _ = strconv.Atoi(s[2])
		f, _ = strconv.Atoi(s[3])
		id, _ = strconv.Atoi(s[4])
		newAction = Action{s[0],s[1],i,f*9,id}
		trace = append(trace, newAction)
		num += 1
	}
	return trace, num
}
