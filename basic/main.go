package main

import (
	"fmt"
)

func main() {
	// 1 Hello World
	fmt.Println("Hello Go")

	// 2 Simple declaration
	var i int

	fmt.Printf("i = %d \n", i)

	// 3 Type inference
	j := 33
	fmt.Printf("j = %d \n", j)

	// 4 Declaring maps
	m := make(map[string]string)
	m["key"] = "value"

	val, ok := m["key"]
	fmt.Printf("Value found in the map: %v %v\n", ok, val)

	// 5 Arrays
	var a [10]int

	fmt.Println("Printing array")
	for i, val := range a {
		fmt.Printf("index = %d value = %d \n", i, val)
	}

	// 6 Slices
	slice := make([]int, 3)

	slice[0] = 1
	slice[1] = 2
	slice[2] = 3

	fmt.Printf("Printing slice: %v \n", slice)

	// 7 structs

	// struct with zero values
	c1 := Complex{}
	fmt.Printf("Struct with zero values: %v \n", c1)

	// stuct initialization
	c2 := Complex{
		Im: 10,
		Re: 10,
	}

	fmt.Printf("Invoking Modulus method on Complex struct: %v \n", c2.Modulus())

	// 7 method receiver on a regular type
	myInt := MyInt(11)

	fmt.Printf("MyInt Methods invocation. IsEven: %v, IsOdd: %v, isPrime: %v \n", myInt.IsEven(), myInt.IsOdd(), myInt.IsPrime())

	// 8 Go routines
	go_routines_example()
}
