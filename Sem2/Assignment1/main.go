package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"
)

type Parameters struct {
	DiskSize   []int  `json:"diskSize"`
	Init       int    `json:"startingPoint"`
	Direction  string `json:"direction"`
	Latency    string `json:"latency"`
	RequestVal []int  `json:"values"`
}

type Request struct {
	Val  int
	Read bool
}

type err struct {
	val int
	err bool
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func checkAbortConditions(params Parameters) bool {
	if params.DiskSize[0] < params.DiskSize[1] {
		fmt.Printf("ABORT(13):upper (%d) < lower (%d)\n", params.DiskSize[0], params.DiskSize[1])
		return true
	}
	if params.Init < params.DiskSize[1] {
		fmt.Printf("ABORT(12):initial (%d) < lower (%d)\n", params.Init, params.DiskSize[1])
		return true
	}
	if params.Init > params.DiskSize[0] {
		fmt.Printf("ABORT(11):initial (%d) > upper (%d)\n", params.Init, params.DiskSize[0])
		return true
	}
	return false
}

func Concurent_fcfs(init chan Parameters, c chan int, er chan []err, waitgroup *sync.WaitGroup) {
	params := <-init
	errors, THM := fcfs(params, params.RequestVal)

	c <- THM
	if errors != nil {
		er <- errors
	}

	waitgroup.Done()
}

func Concurent_scan(init chan Parameters, c chan int, er chan []err, waitgroup *sync.WaitGroup) {
	params := <-init
	errors, THM := scan(params, params.RequestVal)

	c <- THM
	if errors != nil {
		er <- errors
	}
	waitgroup.Done()
}

func Concurent_cscan(init chan Parameters, c chan int, er chan []err, waitgroup *sync.WaitGroup) {
	params := <-init
	errors, THM := cscan(params, params.RequestVal)

	c <- THM
	if errors != nil {
		er <- errors
	}
	waitgroup.Done()
}

func Concurent_look(init chan Parameters, c chan int, er chan []err, waitgroup *sync.WaitGroup) {
	params := <-init
	errors, THM := look(params, params.RequestVal)

	c <- THM
	if errors != nil {
		er <- errors
	}
	waitgroup.Done()
}
func Concurent_clook(init chan Parameters, c chan int, er chan []err, waitgroup *sync.WaitGroup) {
	params := <-init
	errors, THM := clook(params, params.RequestVal)

	c <- THM
	if errors != nil {
		er <- errors
	}
	waitgroup.Done()
}
func Concurent_sstf(init chan Parameters, init_req chan []Request, c chan int, er chan []err, waitgroup *sync.WaitGroup) {
	params := <-init
	reqs := <-init_req
	errors, THM := sstf(params, reqs)

	c <- THM
	if errors != nil {
		er <- errors
	}
	waitgroup.Done()
}
func parseFile(filePath string) (Parameters, []Request) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var params Parameters
	json.Unmarshal(byteValue, &params)
	requests := make([]Request, len(params.RequestVal))
	for index, element := range params.RequestVal {
		requests[index].Val = element
		requests[index].Read = false
	}
	return params, requests
}

func fcfs(params Parameters, requests []int) ([]err, int) {
	var THM = 0
	var currentPos = params.Init
	errors := make([]err, 0)

	for i := 0; i < len(requests); i++ {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	return errors, THM
}
func getShortestIndex(requests []Request, position int, max int) int {
	if len(requests) < 1 {
		return -1
	}
	var diff = max + 1
	var index = -1
	for i := 0; i < len(requests); i++ {
		if requests[i].Read == true {
			continue
		}
		if Abs(requests[i].Val-position) < diff {
			diff = Abs(requests[i].Val - position)
			index = i
		}
	}
	return index
}
func sstf(params Parameters, requests []Request) ([]err, int) {
	var THM = 0
	var currentPos = params.Init
	errors := make([]err, 0)
	var nextIndex = getShortestIndex(requests, currentPos, params.DiskSize[0])
	for i := 0; i < len(requests); i++ {
		if requests[nextIndex].Val > params.DiskSize[0] || requests[nextIndex].Val < params.DiskSize[1] {
			errors = append(errors, err{requests[nextIndex].Val, true})
			requests[nextIndex].Read = true
			continue
		}
		THM += Abs(currentPos - requests[nextIndex].Val)
		currentPos = requests[nextIndex].Val
		requests[nextIndex].Read = true
		errors = append(errors, err{requests[nextIndex].Val, false})
		nextIndex = getShortestIndex(requests, currentPos, params.DiskSize[0])
	}
	return errors, THM
}
func scan(params Parameters, requests []int) ([]err, int) {
	var THM = 0
	var currentPos = params.Init
	var startIndex = len(requests) - 1
	errors := make([]err, 0)
	sort.Ints(requests)
	for i := 0; i < len(requests); i++ {
		if requests[i] > currentPos {
			startIndex = i
			break
		}
	}
	for i := startIndex; i < len(requests); i++ {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	if startIndex < 1 {
		return errors, THM
	}
	THM += Abs(params.DiskSize[0] - currentPos)
	currentPos = requests[startIndex-1]
	THM += Abs(params.DiskSize[0] - currentPos)
	for i := startIndex - 1; i >= 0; i-- {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	return errors, THM
}
func cscan(params Parameters, requests []int) ([]err, int) {
	var THM = 0
	var currentPos = params.Init
	var startIndex = len(requests) - 1
	errors := make([]err, 0)
	sort.Ints(requests)
	for i := 0; i < len(requests); i++ {
		if requests[i] > currentPos {
			startIndex = i
			break
		}
	}
	if startIndex == len(requests)-1 || startIndex == 0 {
		if startIndex == len(requests)-1 {
			THM += Abs(params.DiskSize[0] - params.DiskSize[1])
		}
		for i := 0; i < len(requests); i++ {
			if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
				errors = append(errors, err{requests[i], true})
				continue
			}
			errors = append(errors, err{requests[i], false})
			THM += Abs(currentPos - requests[i])
			currentPos = requests[i]
		}
		return errors, THM
	}
	for i := startIndex; i < len(requests); i++ {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	THM += Abs(params.DiskSize[0] - currentPos)
	currentPos = requests[0]
	THM += Abs(params.DiskSize[0])
	THM += Abs(currentPos)
	for i := 0; i < startIndex; i++ {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	return errors, THM
}
func look(params Parameters, requests []int) ([]err, int) {
	var THM = 0
	var currentPos = params.Init
	var startIndex = len(requests) - 1
	errors := make([]err, 0)
	sort.Ints(requests)
	for i := 0; i < len(requests); i++ {
		if requests[i] > currentPos {
			startIndex = i
			break
		}
	}
	for i := startIndex; i < len(requests); i++ {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	if startIndex < 1 {
		return errors, THM
	}
	THM += Abs(currentPos - requests[startIndex-1])
	currentPos = requests[startIndex-1]
	for i := startIndex - 1; i >= 0; i-- {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	return errors, THM
}
func clook(params Parameters, requests []int) ([]err, int) {
	var THM = 0
	var currentPos = params.Init
	var startIndex = len(requests) - 1
	errors := make([]err, 0)
	sort.Ints(requests)
	for i := 0; i < len(requests); i++ {
		if requests[i] > currentPos {
			startIndex = i
			break
		}
	}
	if startIndex == len(requests)-1 || startIndex == 0 {
		if startIndex == len(requests)-1 {
			THM += Abs(requests[len(requests)-1] - requests[0])
		}
		for i := 0; i < len(requests); i++ {
			if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
				errors = append(errors, err{requests[i], true})
				continue
			}
			errors = append(errors, err{requests[i], false})
			THM += Abs(currentPos - requests[i])
			currentPos = requests[i]
		}
		return errors, THM
	}
	for i := startIndex; i < len(requests); i++ {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	THM += Abs(currentPos - requests[0])
	currentPos = requests[0]
	for i := 0; i < startIndex; i++ {
		if requests[i] > params.DiskSize[0] || requests[i] < params.DiskSize[1] {
			errors = append(errors, err{requests[i], true})
			continue
		}
		errors = append(errors, err{requests[i], false})
		THM += Abs(currentPos - requests[i])
		currentPos = requests[i]
	}
	return errors, THM
}
func printer(prefix string, c chan int, err chan []err) {
	for {
		msg := <-c
		fmt.Println(prefix, msg)
		erro := <-err
		if erro != nil {
			for j := 0; j < len(erro); j++ {
				if erro[j].err == true {
					fmt.Printf("ERROR(15):Request out of bounds: req > upper or < lower\n")
				}
			}
			for k := 0; k < len(erro); k++ {
				if erro[k].err == false {
					fmt.Printf("Servicing Req #%5d\n", erro[k].val)
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func main() {
	var waitgroup sync.WaitGroup
	waitgroup.Add(6)
	inputFile := os.Args[1]
	params, requests := parseFile(inputFile)
	if checkAbortConditions(params) == true {
		fmt.Println("aborted")
		return
	}
	for i := 0; i < len(requests); i++ {
		fmt.Printf("\t\t Requests %5d\n", requests[i].Val)
	}

	resp_fcfs := make(chan int)
	err_fcfs := make(chan []err)

	resp_scan := make(chan int)
	err_scan := make(chan []err)

	resp_cscan := make(chan int)
	err_cscan := make(chan []err)

	resp_look := make(chan int)
	err_look := make(chan []err)

	resp_clook := make(chan int)
	err_clook := make(chan []err)

	resp_sstf := make(chan int)
	err_sstf := make(chan []err)

	init := make(chan Parameters, 6)
	init_req := make(chan []Request, 1)
	init_req <- requests

	for i := 0; i < 6; i++ {
		init <- params
	}

	go Concurent_fcfs(init, resp_fcfs, err_fcfs, &waitgroup)
	go Concurent_scan(init, resp_scan, err_scan, &waitgroup)
	go Concurent_cscan(init, resp_cscan, err_cscan, &waitgroup)
	go Concurent_look(init, resp_look, err_look, &waitgroup)
	go Concurent_clook(init, resp_clook, err_clook, &waitgroup)
	go Concurent_sstf(init, init_req, resp_sstf, err_sstf, &waitgroup)

	go printer("fcfs THM:", resp_fcfs, err_fcfs)
	go printer("scan THM:", resp_scan, err_scan)
	go printer("cscan THM:", resp_cscan, err_cscan)
	go printer("look THM:", resp_look, err_look)
	go printer("clook THM:", resp_clook, err_clook)
	go printer("sstf THM:", resp_sstf, err_sstf)

	time.Sleep(1 * time.Second)

}
