package main

import (
	"fmt"
	"math"
)

var numArr []int
var helperArr []int

func main() {
	result, _ := Sort([]int{1, 43, 5, 34, 2, 5, 6, 6, 10, 21}, 3)
	fmt.Printf("%+v\n", result)
}

func Sort(arr []int, k int) (result []int, err error) {
	numArr = arr
	helperArr = make([]int, len(arr))
	copy(helperArr, arr)

	partitioning(k, 0, len(helperArr)-1)

	return arr, nil
}

func getSizes(k, start, end int) (totalSize, sizePartitions, sizeEndPartition int) {
	totalSize = end - start + 1
	sizePartitions = int(math.Max(float64(totalSize/k), 1.0))
	sizeEndPartition = totalSize - (sizePartitions * (k - 1))
	return
}

func partitioning(k, start, end int) {
	_, sizePartitions, sizeEndPartition := getSizes(k, start, end)

	// Recurse until the size of each partition is 1
	if sizePartitions > 1 {
		for i := 0; i < k; i++ {

			var newEndPart int
			if i == k-1 {
				newEndPart = i*sizePartitions + start + sizeEndPartition - 1
			} else {
				newEndPart = i*sizePartitions + start + sizePartitions - 1
			}

			// Recurse for each partition
			partitioning(k, i*sizePartitions+start, newEndPart)
		}

	} else {
		// if the size of the partitions is 1
		if sizeEndPartition > 1 {
			partitioning(k, end-sizeEndPartition, end)
		}
	}

	// Once we are done recursing, start merging
	merge(k, start, end)
}

func merge(k, low, high int) {
	totalSize, sizePartitions, sizeLastPartition := getSizes(k, low, high)

	if totalSize < k {
		merge(totalSize, low, high)
		return
	}

	var (
		indices                           = make([]int, k)
		count                         int = low
		min, minPosition, currentPart int
	)

	for count <= high {

		/** Iterate through first element in every partition and find the minimum */
		min = math.MaxInt32
		minPosition = 0
		currentPart = low
		for i := 0; i < k; i++ {

			if i == k-1 && indices[i] == sizeLastPartition {
				break
			}

			if indices[i] == sizePartitions && i != k-1 {
				currentPart += sizePartitions
				continue
			}

			if helperArr[currentPart+indices[i]] < min {
				min = helperArr[currentPart+(indices[i])]
				minPosition = i
			}

			currentPart += sizePartitions
		}

		numArr[count] = min
		count++
		indices[minPosition]++
	}

	for i := low; i <= high; i++ {
		helperArr[i] = numArr[i]
	}
}
