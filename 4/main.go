package main

import (
	"fmt"
)

func solveA(start, end int) int {

	sum := 0
	for a := start / 100000; a <= end/100000; a++ {
		for b := a; b < 10; b++ {
			for c := b; c < 10; c++ {
				for d := c; d < 10; d++ {
					for e := d; e < 10; e++ {
						for f := e; f < 10; f++ {
							consecutive := a == b || b == c || c == d || d == e || e == f
							if !consecutive {
								continue
							}
							val := (1e5 * a) + (1e4 * b) + (1e3 * c) + (1e2 * d) + (1e1 * e) + (f)
							if val >= start && val <= end {
								sum++
							}
						}
					}
				}
			}
		}
	}
	return sum

}

func solveB(start, end int) int {

	sum := 0
	for a := start / 100000; a <= end/100000; a++ {
		for b := a; b < 10; b++ {
			for c := b; c < 10; c++ {
				for d := c; d < 10; d++ {
					for e := d; e < 10; e++ {
						for f := e; f < 10; f++ {

							// Find a run of exactly length 2
							run := 1
							arr := []int{a, b, c, d, e, f}
							consecutive := false
							for x := 1; x < len(arr); x++ {
								if arr[x] == arr[x-1] {
									run++
								} else {
									if run == 2 {
										consecutive = true
									}
									run = 1
								}
							}
							if run == 2 {
								consecutive = true
							}

							if !consecutive {
								continue
							}
							val := (1e5 * a) + (1e4 * b) + (1e3 * c) + (1e2 * d) + (1e1 * e) + (f)
							if val >= start && val <= end {
								sum++
							}
						}
					}
				}
			}
		}
	}
	return sum
}

func main() {
	fmt.Printf("A: %d\n", solveA(125730, 579381))
	fmt.Printf("B: %d\n", solveB(125730, 579381))
}
