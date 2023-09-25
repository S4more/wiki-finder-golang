package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"slices"
	"sync"
	"testing"

	file "github.com/s4more/go-wiki/loader"
	utils "github.com/s4more/go-wiki/utils"
)

func TestMain(m *testing.M) {
	f, err := os.Create("profile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	defer utils.Timer("main")()
	println("Starts reading file.")
	bytes := file.ReadFile()
	println("Finished reading file.")
	outgoingLinks = file.UnmarshalMessage(bytes)
	println("Finished json parse.")

	i := 1

	defer utils.Timer("0-50")()
	for i < 1000 {
		var wg sync.WaitGroup
		isClosed = false
		visited = make([]uint32, 100)
		nodeCh := make(chan uint32)
		stopCh := make(chan struct{})

		go findOutgoingLink(0, uint32(i), 0, nodeCh, stopCh, &wg)

		node := <-nodeCh
		fmt.Printf("Received %v on path %v", node, i)
		isClosed = true
		close(stopCh)
		wg.Wait()

		if node != MAX_VAL {
			fmt.Printf("Found response %v -> %v \n", i, node)
			if !slices.Contains(outgoingLinks[node], uint32(i)) {
				fmt.Printf("node %v does not have an outgoing link to %v", node, i)
			}
		} else {
			fmt.Printf("No response.")
		}
		i += 1
	}
}
