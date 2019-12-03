package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	wires := make([][]string, 0)
	for sc.Scan() {
		wire := strings.Split(sc.Text(), ",")
		wires = append(wires, wire)
	}

	// map of coordinate pair to which wires touch it
	m := make(map[string][]int)
	for idx, wire := range wires {

		var x, y int
		for _, c := range wire {
			xd, yd := 0, 0
			switch c[0] {
			case 'R':
				xd = 1
			case 'L':
				xd = -1
			case 'U':
				yd = 1
			case 'D':
				yd = -1
			default:
				panic("invalid cmd")
			}

			for i := mustParse(c[1:]); i > 0; i-- {
				x += xd
				y += yd
				key := fmt.Sprintf("%d:%d", x, y)
				m[key] = append(m[key], idx)
			}
		}
	}

	// pick best of intersections
	best := 9999999999
	for point, wires := range m {
		if len(wires) > 1 {
			// make sure at least 2 unique wires
			unique := make(map[int]bool)
			for _, i := range wires {
				unique[i] = true
			}
			if len(unique) < 2 {
				continue
			}

			pair := strings.Split(point, ":")
			dist := abs(mustParse(pair[0])) + abs(mustParse(pair[1]))
			if dist < best {
				best = dist
			}

		}
	}

	return best
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func mustParse(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func solveB(r io.Reader) int {
	return -1
}

func main() {
	input := open("input.txt")
	fmt.Printf("A: %d\n", solveA(input))
	//input = open("input.txt")
	//fmt.Printf("B: %d\n", solveB(input))
}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}
