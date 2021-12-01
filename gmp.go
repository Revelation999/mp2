package main

import (
	"runtime"
)

var puzzleSolved = make(chan bool)

// MaxParallelism returns the max GOMAXPROCS value, by comparing runtime.GOMAXPROCS
// with number of CPU threads, and returning larger value
func MaxParallelism() int {
	maxProcess := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcess < numCPU {
		return maxProcess
	}
	return numCPU
}
