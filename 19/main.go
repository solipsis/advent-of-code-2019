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

	c := newCPU(program)

	size := 99

	row := 15
	col := 10
	colStart := col
	prevRes := 1

	next := []int{col, row}
	c.input = append(c.input, next...)

	//62817 too low
	for {

		out, err := c.run()
		switch err {
		case nil:
			c = newCPU(program)
			c.input = append(c.input, next...)
			continue
		case errSuspend:
			fmt.Printf("r: %d, c: %d, v: %d\n", row, col, out)
			if out == 1 {

				if prevRes == 0 {
					colStart = col
				}

				// check size x size grid
				// right
				c = newCPU(program)
				c.input = append(c.input, []int{col + size, row}...)
				r, _ := c.run()
				// down
				c = newCPU(program)
				c.input = append(c.input, []int{col, row + size}...)
				d, _ := c.run()
				// right down
				c = newCPU(program)
				c.input = append(c.input, []int{col + size, row + size}...)
				rd, _ := c.run()

				// we found the solution
				if r == 1 && d == 1 && rd == 1 {
					return (col * 10000) + row
				}

				// check next column
				col++
				next = []int{col, row}
			} else {
				// if 2 consecutive misses go to new row
				if prevRes == 0 {
					row++
					col = colStart
					prevRes = 1
					next = []int{col, row}
					continue
				} else {
					// try next column
					col++
					next = []int{col, row}
				}
			}
			prevRes = out

		case errNeedInput:
			panic("shouldn't need more input?")

		default:
			panic(err)
		}
	}

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
	//fmt.Printf("A: %d\n", solveA(input))
	//input = open("input.txt")
	fmt.Printf("B: %d\n", solveB(input))
}
