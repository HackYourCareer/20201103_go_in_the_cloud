package main

import "fmt"

func main() {

	// 1) Scan ports
	scanPortsExample()

	// 2) Run pipeline
	//runPipelinesExample()

	// 3) Run pipeline with fan-in-fan-out pattern
	runPipelinesFanInFanOutExample()
}

func scanPortsExample() {
	openPorts := scanPorts("httpbin.org", 1, 1024)
	fmt.Printf("Open ports : %v /n", openPorts)
}

func scanPortsConcurrentlyExample() {
	openPorts := scanPortsConcurrently("httpbin.org", 1, 1024)
	fmt.Printf("Open ports : %v /n", openPorts)
}

func runPipelinesExample() {
	primes := primesGenerator(1, 2000)
	primesTransformed := add(square(primes), 7)

	readFromPipeline(primesTransformed)
}

func runPipelinesFanInFanOutExample() {
	primes1 := primesGenerator(1, 500)
	primes2 := primesGenerator(501, 1000)
	primes3 := primesGenerator(1001, 1501)
	primes4 := primesGenerator(1501, 2000)

	primes := merge(primes1, primes2, primes3, primes4)
	primesTransformed := add(square(primes), 7)

	readFromPipeline(primesTransformed)
}
