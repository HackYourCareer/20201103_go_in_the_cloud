package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func scanPorts(address string, from, to int) []int {
	ports := make([]int, 0)
	for port := from; port <= to; port++ {
		fmt.Printf("Scanning port: %d \n", port)
		a := fmt.Sprintf("%s:%d", address, port)

		conn, err := net.DialTimeout("tcp", a, time.Second*5)

		if err == nil {
			conn.Close()
			ports = append(ports, port)
		} else {
			fmt.Printf("Failed to open port %d: %s \n", port, err)
		}
	}

	return ports
}

func scanPortsConcurrently(address string, from, to int) []int {
	results := make(chan int)

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(to - from + 1)

		for i := from; i <= to; i++ {
			startWorker(address, i, results, &wg)
		}
		wg.Wait()
		close(results)
	}()

	return readFromChannel(results)
}

func startWorker(address string, port int, results chan<- int, wg *sync.WaitGroup) {
	go func() {
		fmt.Printf("Scanning port: %d \n", port)
		a := fmt.Sprintf("%s:%d", address, port)

		conn, err := net.DialTimeout("tcp", a, time.Second*5)

		if err == nil {
			conn.Close()
			results <- port
		} else {
			fmt.Printf("Failed to open port %d: %s \n", port, err)
		}
		wg.Done()
	}()
}

func readFromChannel(results <-chan int) []int {
	ports := make([]int, 0)

	for port := range results {
		ports = append(ports, port)
		fmt.Printf("Value %d read from channel \n", port)
	}

	return ports
}
