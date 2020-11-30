package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	//"time"

	//"time"

	//"runtime"
	"strings"
)

const (
	iterationsNum = 2
	goroutinesNum = 2
	quotaLimit    = 2
)

/*
func main() {
	wg := &sync.WaitGroup{}
	quotaCh := make(chan struct{}, quotaLimit) // ratelim.go
	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		go startWorker(i, wg, quotaCh)
	}
	time.Sleep(time.Millisecond)
	wg.Wait()
}
*/
func startWorker(in int, wg *sync.WaitGroup, quotaCh chan struct{}) {
	quotaCh <- struct{}{} // ratelim.go, берём свободный слот
	defer wg.Done()
	for j := 0; j < iterationsNum; j++ {

		//fmt.Printf(formatWork(in, j))
		select {}
	}
	<-quotaCh // ratelim.go, возвращаем слот
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat(" ", in), "*",
		strings.Repeat(" ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("*", j))
}

func main() {
	var ok = true
	var received uint32
	jobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			fmt.Println("job1-out")
			//fmt.Println(in)
			time.Sleep(10 * time.Millisecond)
			currRecieved := atomic.LoadUint32(&received)
			fmt.Println("job1-received")
			fmt.Println(received)
			// в чем тут суть
			// если вы накапливаете значения, то пока вся функция не отрабоатет - дальше они не пойдут
			// тут я проверяю, что счетчик увеличился в следующей функции
			// это значит что туда дошло значение прежде чем текущая функция отработала
			if currRecieved == 0 {
				ok = false
			}

		}),

		job(func(in, out chan interface{}) {
			fmt.Println("job2-in")
			//fmt.Println(in)
			for _ = range in {
				atomic.AddUint32(&received, 1)
				fmt.Println("job2-received")
				fmt.Println(received)
			}
		}),
	}
	//start := time.Now()

	ExecutePipeline(jobs...)
	//fmt.Scanln()
	//end := time.Since(start)
}

func ExecutePipeline(job ...job) {
	in := make(chan interface{}, 5)
	out := make(chan interface{}, 5)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(in, out chan interface{}) {
		for _, job := range job {
			job(in, out)
			fmt.Println("out")
			fmt.Println(out)
			in <- out
			fmt.Println("in2")
			fmt.Println(in)

		}
		//fmt.Println("in")
		//fmt.Println(in)

		wg.Done()
	}(in, out)

	/*
		for _, job := range job {
			fmt.Println("durr1")
			job(in, out)
			in <- out



		}
	*/
	/*
		for _, job1 := range job {
			wg.Add(1)
			go func(func(in, out chan interface{})) {
				//fmt.Println("in")
				//fmt.Println(in)
				job1(in, out)
				fmt.Println("out")
				fmt.Println(out)
				in <- out
				fmt.Println("in2")
				fmt.Println(in)
				wg.Done()
			}(job1)
		}
	*/

	wg.Wait()

	/*

		for _, job := range job {
			go func() {
				for _, n := range nums {
					out <- n
				}
				close(out)
			}()
			go job(in, out)

			change <-out

			in <- change
			//fmt.Print( )
			//SingleHash(in, out)
			//data, ok := dataRaw.(string)
			//println()out
		}
	*/

	/*
		select {
		case val := <-in:
			fmt.Println("ch1 val", val)
		case out <- 1:
			fmt.Println("put val to out")
		default:
			fmt.Println("default case")
		}

		//runtime.Goshed()
		go func() {
			for _, job := range job {
				job(in, out)
				//fmt.Print( )
				SingleHash(in, out)
				//data, ok := dataRaw.(string)
				//println()out
			}
			close(out)
		}()



	*/

	/*

		return out


		for i, job := range job {

			select {
			case joba
			case <-ctx.Done():
				break LOOP
			case foundBy := <-result:
				totalFound++
				fmt.Println("result found by", foundBy)
			}
		}
	*/
	//time.Sleep(time.Millisecond)
	//wg.Wait() // wait_2.go ожидаем, пока waiter.Done() не приведёт счетчик к 0
}

//var SingleHash job

var SingleHash = func(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		panic("can't convert result data to string")
	}
	fmt.Println(data)
	out <- DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
}

var MultiHash = func(in, out chan interface{}) {
	/*
		dataRaw := <- in
		data, ok := dataRaw.(string)
		if !ok {
			panic("can't convert result data to string")
		}
		out <- DataSignerCrc32(in.(string) + data) + "~" + DataSignerCrc32(DataSignerMd5(data))
		return fmt.Sprintln(strings.Repeat(" ", in), "*",
			strings.Repeat(" ", goroutinesNum-in),
			"th", in,
			"iter", j, strings.Repeat("*", j))
	*/
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		panic("cant convert result data to string")
	}
	out <- DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
}

var CombineResults = func(in, out chan interface{}) {

}

//(in, out chan interface{}) {

/*
func SingleHash(data string) string {
	return DataSignerCrc32(data)
}
*/

/*

func job(f func) string {
	return fmt.Sprintln(strings.Repeat(" ", in), "*",
		strings.Repeat(" ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("*", j))
}
*/
