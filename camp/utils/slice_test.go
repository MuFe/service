package utils

import (
	"fmt"
	"testing"
)

func TestIntersectInt32(t *testing.T) {
	result := IntersectInt32([]int32{1, 2, 3}, []int32{2, 3, 4, 5})
	fmt.Println("result", result)
}
