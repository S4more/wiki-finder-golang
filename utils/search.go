package utils

import (
	"slices"
)

func SimpleSearch(arr []uint32, key uint32) bool {
	for _, v := range arr {
		if v == key {
			return true
		}
	}

	return false
}

func BinarySearch(arr []uint32, key uint32) bool {
	_, exists := slices.BinarySearch(arr, key)
	return exists
}
