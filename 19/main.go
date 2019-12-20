package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var grid [][]int

func solveB(r io.Reader) int {
	return -1
}

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	program := make([]int, 0)

	for sc.Scan() {
		nums := strings.Split(sc.Text(), ",")
		for _, v := range nums {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			program = append(program, i)
		}
	}

	sum := 0

	c := newCPU(program)
	coords := []int{}
	// construct input
	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			coords = append(coords, []int{x, y}...)
		}
	}

	for len(coords) > 0 {

		out, err := c.run()
		switch err {
		case nil:
			c = newCPU(program)
			continue
		case errSuspend:
			fmt.Println("OUTPUT: ", out)
			if out == 1 {
				sum++
			}

		case errNeedInput:
			fmt.Println("providing input", coords[0])
			c.input = append(c.input, coords[0])
			coords = coords[1:]

		default:
			panic(err)
		}
	}

	return sum
}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	input := open("input.txt")
	fmt.Printf("A: %d\n", solveA(input))
	//input = open("input.txt")
	//fmt.Printf("B: %d\n", solveB(input))
}
