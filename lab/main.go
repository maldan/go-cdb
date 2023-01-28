package main

import (
	"fmt"
	"time"
)

func Find[T comparable](arr []T, val T) []T {
	out := make([]T, 0)
	for i := 0; i < len(arr); i++ {
		if arr[i] == val {
			out = append(out, arr[i])
		}
	}
	return out
}

func main() {
	// fmt.Printf("%v", "sas")
	/*ss := map[int]any{}
	for i := 0; i < 1_000_000; i++ {
		ss[i] = i
	}
	fmt.Printf("%v", ss[100000])*/

	ss := make([]string, 1_000_000)
	for i := 0; i < 1_000_000; i++ {
		ss[i] = fmt.Sprintf("%v", i)
	}

	// Find
	t := time.Now()
	x := Find(ss, "128000")
	fmt.Printf("%v - %v\n", x, time.Since(t))

	time.Sleep(time.Hour)
}
