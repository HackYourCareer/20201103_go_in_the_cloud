package main

import (
	"fmt"
	"sync"
)

func primesGenerator(from, to int) <-chan int {
	out := make(chan int)

	go func() {
		for i := from; i <= to; i++ {
			if MyInt(i).IsPrime() {
				out <- i
			}
		}

		close(out)
	}()

	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range in {
			out <- v * v
		}

		close(out)
	}()

	return out
}

func add(in <-chan int, valueToAdd int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range in {
			out <- v + valueToAdd
		}

		close(out)
	}()

	return out
}

func readFromPipeline(in <-chan int) {
	for port := range in {
		fmt.Printf("Received from pileline: %d \n", port)
	}
}

func merge(channels ...<-chan int) <-chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}

	wg.Add(len(channels))

	readAndWriteFunc := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			out <- v
		}
	}

	go func() {
		for _, ch := range channels {
			go readAndWriteFunc(ch)
		}
		wg.Wait()
		close(out)
	}()

	return out
}
