package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func go_routines_example() {
	wg := sync.WaitGroup{}

	wg.Add(2)
	val := 0

	pval := &val

	go modify(pval, &wg)
	go modify(pval, &wg)

	wg.Wait()
	fmt.Println(val)
}

func go_routines_example_atomic() {
	wg := sync.WaitGroup{}

	wg.Add(2)
	val := int64(0)

	pval := &val

	go modify_atomic(pval, &wg)
	go modify_atomic(pval, &wg)

	wg.Wait()
	fmt.Println(val)
}

func go_routines_example_mutex() {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	wg.Add(2)
	val := int64(0)

	pval := &val

	go modify_mutex(pval, &wg, &mutex)
	go modify_mutex(pval, &wg, &mutex)

	wg.Wait()
	fmt.Println(val)
}

func modify(val *int, wg *sync.WaitGroup) {
	for i := 0; i < 10000; i++ {
		*val = *val + 1
	}
	wg.Done()
}

func modify_atomic(val *int64, wg *sync.WaitGroup) {
	for i := 0; i < 100000; i++ {
		atomic.AddInt64(val, 1)
	}
	wg.Done()
}

func modify_mutex(val *int64, wg *sync.WaitGroup, mutex *sync.Mutex) {
	for i := 0; i < 100000; i++ {
		mutex.Lock()
		*val++
		mutex.Unlock()
	}
	wg.Done()
}
