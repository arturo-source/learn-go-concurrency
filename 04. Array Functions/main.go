package main

import (
	"runtime"
	"sync"
)

func Map[T1, T2 any](arr []T1, f func(item T1, index int) T2) []T2 {
	arrT2 := make([]T2, len(arr))

	numworkers := runtime.NumCPU()

	var wg sync.WaitGroup
	wg.Add(numworkers)

	worker := func(startIndex, endIndex int) {
		defer wg.Done()

		for i := startIndex; i < endIndex; i++ {
			t2 := f(arr[i], i)
			arrT2[i] = t2
		}
	}

	chunkSize := len(arr) / numworkers
	for i := 0; i < numworkers; i++ {
		startIndex := i * chunkSize
		endIndex := (i + 1) * chunkSize
		if i == numworkers-1 {
			endIndex = len(arr)
		}

		go worker(startIndex, endIndex)
	}

	wg.Wait()

	return arrT2
}

func Filter[T any](arr []T, f func(item T, index int) bool) []T {
	arrBooled := make([]bool, len(arr))

	numworkers := runtime.NumCPU()

	var wg sync.WaitGroup
	wg.Add(numworkers)

	worker := func(startIndex, endIndex int) {
		defer wg.Done()

		for i := startIndex; i < endIndex; i++ {
			arrBooled[i] = f(arr[i], i)
		}
	}

	chunkSize := len(arr) / numworkers
	for i := 0; i < numworkers; i++ {
		startIndex := i * chunkSize
		endIndex := (i + 1) * chunkSize
		if i == numworkers-1 {
			endIndex = len(arr)
		}

		go worker(startIndex, endIndex)
	}

	wg.Wait()

	arrFiltered := make([]T, 0)
	for i, val := range arrBooled {
		if val {
			arrFiltered = append(arrFiltered, arr[i])
		}
	}

	return arrFiltered
}

type Number struct {
	a int
}

func main() {
	const n = 100000
	arr := make([]Number, 0, n)
	for i := range n {
		arr = append(arr, Number{i})
	}

	_ = Map(arr, func(t Number, index int) int { return t.a })
	_ = Filter(arr, func(item Number, index int) bool { return item.a > 100 })
}
