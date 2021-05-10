package main

import (
	"fmt"
)

func main() {
	// Создаем необходимые каналы.
	c := gen(2, 3, 6)
	v := sq(c)
	out := sq2(v)
	// Выводим значения.
	fmt.Println(<-out) // 4
	fmt.Println(<-out) // 9
	fmt.Println(<-out) // 9
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func sq2(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}
