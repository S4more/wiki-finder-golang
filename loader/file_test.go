package file_test

import (
	"testing"
)

var neighbours []uint32
var rang int = 2000

func TestEytzinger(t *testing.T) {
	neighbours = make([]uint32, rang)
	for i := 0; i < rang; i++ {
		neighbours[i] = uint32(i)
	}

	SyncBinarySearch(278, 0, len(neighbours))
	// eytzinger.Sort(neighbours)
	// eytzinger.Search(neighbours, 278)

}

func SyncBinarySearch(x uint32, start int, end int) bool {
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
