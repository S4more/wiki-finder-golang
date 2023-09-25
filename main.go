package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	_ "net/http/pprof"

	file "github.com/s4more/go-wiki/loader"
	utils "github.com/s4more/go-wiki/utils"
	"github.com/tevino/abool/v2"
)

var outgoingLinks [][]uint32

const MAX_VAL = 4294967295
const MAX_DEPTH = 4

func BinarySearch(a uint32, x uint32) bool {
	start := 0
	end := len(outgoingLinks[a]) - 1
	for start <= end {
		mid := (start + end) / 2
		if outgoingLinks[a][mid] == x {
			return true
		} else if outgoingLinks[a][mid] < x {
			start = mid + 1
		} else if outgoingLinks[a][mid] > x {
			end = mid - 1
		}
	}
	return false
}

func hasNeighbour(startLink, goal uint32) bool {
	// outgoingNeighbours := outgoingLinks[startLink]
	return BinarySearch(startLink, goal)
}

func findOutgoingLink(startLink uint32, goal uint32, currentDepth uint32, returnCh chan<- uint32, stopCh <-chan struct{}, isClosed *abool.AtomicBool) {
	if isClosed.IsSet() {
		return
	}
	send := func(data uint32) {
		select {
		case <-stopCh:
			return
		case returnCh <- data:
		}
	}

	// visited = append(visited, startLink)
	if hasNeighbour(startLink, goal) {
		send(startLink)
		return
	}

	if currentDepth == MAX_DEPTH {
		return
	}

	for _, l := range outgoingLinks[startLink] {
		if isClosed.IsSet() {
			return
		}

		// if currentDepth+2 == MAX_DEPTH {
		// 	go findOutgoingLink(l, goal, currentDepth+1, returnCh, stopCh, isClosed)
		// } else {
		findOutgoingLink(l, goal, currentDepth+1, returnCh, stopCh, isClosed)
		// }

	}
}

func main() {
	defer utils.Timer("main")()
	println("Starts reading file.")
	bytes := file.ReadFile()
	println("Finished reading file.")
	outgoingLinks = file.UnmarshalMessage(bytes)
	println("Finished json parse.")

	i := 1

	f, _ := os.Create("cpu_profile.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	_ = http.ListenAndServe("0.0.0.0:8081", nil)
	// 	wg.Done()
	// }()

	defer utils.Timer("0-50")()
	for i < 1000 {
		cond := abool.New()
		nodeCh := make(chan uint32)
		stopCh := make(chan struct{})

		go findOutgoingLink(0, uint32(i), 0, nodeCh, stopCh, cond)

		node := <-nodeCh
		cond.Set()
		close(stopCh)

		if node != MAX_VAL {
			fmt.Printf("Found response %v -> %v \n", i, node)
			// if !slices.Contains(outgoingLinks[node], uint32(i)) {
			// 	fmt.Printf("node %v does not have an outgoing link to %v", node, i)
			// }
		} else {
			fmt.Printf("No response.")
		}
		i += 1
	}

	// wg.Wait()
}
