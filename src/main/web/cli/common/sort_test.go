package common

import (
	"fmt"
	"testing"
)

func TestSort(t *testing.T) {
	nums := []int{4, 5, 21, 1, 34, 6}
	MySort(nums, "asc")
	fmt.Printf("asc sorted : %v", nums)
	MySort(nums, "desc")
	fmt.Printf("desc sorted : %v", nums)
}
