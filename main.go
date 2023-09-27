package main

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"

	_ "net/http/pprof"

	file "github.com/s4more/go-wiki/loader"
	utils "github.com/s4more/go-wiki/utils"
	"golang.org/x/sync/semaphore"
)

var outgoingLinks [][]uint32
var lengths []int

const MAX_VAL = 4294967295
const MAX_DEPTH = 4

func SyncBinarySearch(a uint32, x uint32, start int, end int) bool {
	neighbours := outgoingLinks[a]
	for start <= end {
		mid := (start + end) >> 1
		if neighbours[mid] < x {
			start = mid + 1
		} else if neighbours[mid] > x {
			end = mid - 1
		} else {
			return true
		}
	}
	return false
}

func hasNeighbour(startLink, goal uint32) bool {
	length := lengths[startLink]
	return SyncBinarySearch(startLink, goal, 0, length-1)
}

func findOutgoingLink(startLink uint32, goal uint32, currentDepth uint32, visited map[uint32]struct{}) uint32 {
	_, visitedBefore := visited[startLink]
	if visitedBefore {
		return MAX_VAL
	}

	if currentDepth+1 == MAX_DEPTH {
		visited[startLink] = struct{}{}
	}

	if hasNeighbour(startLink, goal) {
		return startLink
	}

	if currentDepth == MAX_DEPTH {
		return MAX_VAL
	}

	for _, l := range outgoingLinks[startLink] {
		val := findOutgoingLink(l, goal, currentDepth+1, visited)
		if val != MAX_VAL {
			return val
		}
	}

	return MAX_VAL
}

func main() {
	defer utils.Timer("main")()
	println("Starts reading file.")
	bytes := file.ReadFile()
	println("Finished reading file.")
	outgoingLinks = file.UnmarshalMessage(bytes)

	lengths = make([]int, len(outgoingLinks))
	// Precomputing the lengths of the arrays
	for i, val := range outgoingLinks {
		leng := len(val)
		lengths[i] = leng - 1
	}

	println("Finished json parse.")

	i := 0
	f, _ := os.Create("cpu_profile.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	sem := semaphore.NewWeighted(16)
	answers := make([]uint32, 25_000)
	ctx := context.TODO()

	defer utils.Timer("0-50")()
	for range answers {
		if err := sem.Acquire(ctx, 1); err != nil {
			println("Couldn't acquire semaphore")
			println(err)
			break
		}

		go func(i uint32) {
			defer sem.Release(1)
			visited := make(map[uint32]struct{})

			node := findOutgoingLink(0, uint32(i), 0, visited)
			answers[i] = node
		}(uint32(i))

		i += 1
		if i%1000 == 0 {
			fmt.Printf("Progress: %v\n", i)
		}
	}
}
