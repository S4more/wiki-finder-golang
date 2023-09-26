package main_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("Hello World")
	size := 15000

	arr := make([]uint32, size)
	for i := range arr {
		arr[i] = uint32(rand.Intn(60_000_000))
	}

	goal := arr[rand.Intn(size)]
	fmt.Printf("Goal: %v\n", goal)

	ch := make(chan bool, 2)

	start := time.Now()
	val := SyncBinarySearch(arr, goal, 0, size-1)
	elapsed := time.Since(start)
	log.Printf("Elapsed time sync: %v, %v\n", elapsed, val)

	start = time.Now()
	go AsyncBinarySearch(arr, goal, 0, size>>1-1, ch)
	go AsyncBinarySearch(arr, goal, size>>1, size-1, ch)
	go AsyncBinarySearch(arr, goal, size>>1, size-1, ch)

	i := 0
	found := false
	for i < 2 {
		val := <-ch
		if val {
			found = true
			break
		}
		i++
	}
	elapsed = time.Since(start)
	log.Printf("Elapsed time: %v, %v\n", elapsed, found)
}

func SyncBinarySearch(neighbours []uint32, goal uint32, start int, end int) bool {

	for start <= end {
		mid := (start + end) >> 1
		if neighbours[mid] == goal {
			return true
		} else if neighbours[mid] < goal {
			start = mid + 1
		} else if neighbours[mid] > goal {
			end = mid - 1
		}
	}
	return false
}

func AsyncBinarySearch(neighbours []uint32, goal uint32, start int, end int, rch chan bool) {
	for start <= end {
		mid := (start + end) >> 1
		if neighbours[mid] == goal {
			rch <- true
		} else if neighbours[mid] < goal {
			start = mid + 1
		} else if neighbours[mid] > goal {
			end = mid - 1
		}
	}
	rch <- false
}
