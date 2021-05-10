package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

const (
	goroutines = 4
)

var combinedResults []string
var globalIndex int

type log struct {
	id  string
	log string
}

var logs = map[string]*log{}
var mutex = &sync.Mutex{}

func main() {
	jobs := []job{
		job(func(in, out chan interface{}) {
			SingleHash2(in, out)
		}),
		job(func(in, out chan interface{}) {
			MultiHash2(in, out)
		}),
		/*
			job(func(in, out chan interface{}) {
				CombineResults2(in, out)
			}),

		*/
	}

	wg := &sync.WaitGroup{}
	for index := 0; index < goroutines; index++ {
		wg.Add(1)
		go func(index int) {
			globalIndex = index
			defer wg.Done()
			ExecutePipeline2(jobs...)
		}(index)
	}
	wg.Wait()

	// Log print
	var sortedLogs = map[string]string{}
	ids := make([]string, 0, len(logs))

	for k, val := range logs {
		sortedLogs[val.id] = val.log
		ids = append(ids, logs[k].id)
	}
	sort.Strings(ids)

	for _, key := range ids {
		fmt.Println(sortedLogs[key])
	}

	fmt.Println(len(combinedResults))

	sort.Strings(combinedResults)
	for _, val := range combinedResults {
		fmt.Printf("%s_", val)
	}

}

func ExecutePipeline2(jobs ...job) {
	in := make(chan interface{}, 1)
	out := make(chan interface{}, 1)
	//cancelChan := make(chan interface{}, 1)

	in <- strconv.Itoa(globalIndex)
	for _, currentJob := range jobs {
		currentJob(in, out)
		val := <-out
		in <- val
	}
}

var SingleHash2 = func(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		panic("can't convert result data to string")
	}
	const name = "SingleHash"

	var md5 string
	var crc32md5 string
	var crc32 string

	dataStrMd5 := make(chan string, 1)

	mutex.Lock()
	dataStrMd5 <- DataSignerMd5(data)
	mutex.Unlock()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		md5 = <-dataStrMd5
		crc32md5 = DataSignerCrc32(md5)
		close(dataStrMd5)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		crc32 = DataSignerCrc32(data)
	}()

	wg.Wait()

	result := crc32 + "~" + crc32md5

	localLog := fmt.Sprintf("%s %s data %s\n", data, name, data)
	localLog += fmt.Sprintf("%s %s md5(data) %s\n", data, name, md5)
	localLog += fmt.Sprintf("%s %s crc32(md5(data)) %s\n", data, name, crc32md5)
	localLog += fmt.Sprintf("%s %s crc32(data) %s\n", data, name, crc32)
	localLog += fmt.Sprintf("%s %s result %s\n", data, name, result)

	logs[result] = &log{data, localLog}

	out <- result
}

var MultiHash2 = func(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		panic("can't convert result data to string")
	}
	const Offset = 0
	const Limit = 6

	var result string

	wg := &sync.WaitGroup{}

	outArr := [Limit]string{}
	for i := Offset; i < Limit; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			outArr[i] = DataSignerCrc32(strconv.Itoa(i) + data)
		}(i)
	}

	wg.Wait()

	log := logs[data]
	for i := Offset; i < Limit; i++ {
		result += outArr[i]
		logs[data].log += fmt.Sprintf("%s MultiHash: crc32(th+step1)) %d %s\n", data, i, outArr[i])
	}
	log.log += fmt.Sprintf("%s MultiHash result: %s\n", data, result)

	out <- result
}

var CombineResults2 = func(in, out chan interface{}) {
	out <- ""
	if globalIndex+1 == goroutines {
		//return
	}

	for {
		select {
		case val := <-in:
			data, ok := val.(string)
			if !ok {
				panic("can't convert result data to string")
			}
			fmt.Println("data:" + data)
			combinedResults = append(combinedResults, data)
		case <-out:
			//panic("STOP")
			return
		}
	}
}
