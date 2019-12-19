package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type parameterMode int

var grid [][]int

const (
	position parameterMode = iota
	immediate
	relative
)

type cpu struct {
	state []int
	input []int
	pc    int
	rb    int // relative base
}

func newCPU(program []int) *cpu {

	pCpy := make([]int, len(program)+10000)
	copy(pCpy, program)

	return &cpu{
		state: pCpy,
		pc:    0,
	}

}

var errNeedInput = errors.New("Need additional input")
var errSuspend = errors.New("Execution supspended")

func (c *cpu) run() (int, error) {
	//fmt.Println("##########################################")
	//fmt.Printf("input: %v\n", c.input)

	for {

		op := c.state[c.pc]

		opCode := op % 100
		m1 := (op / 100) % 10
		m2 := (op / 1000) % 10
		m3 := (op / 10000) % 10
		var p1, p2, p3 = c.pc + 1, c.pc + 2, c.pc + 3

		//fmt.Printf("op: %d, m1: %d, m2: %d, m3: %d, c.pc: %d\n", op, m1, m2, m3, c.pc)
		//fmt.Printf("%v\n", c.state[c.pc:c.pc+4])
		switch opCode {
		// Add
		case 1:
			c.set(p3, m3, c.fetch(p1, m1)+c.fetch(p2, m2))
			c.pc += 4
		// Mul
		case 2:
			c.set(p3, m3, c.fetch(p1, m1)*c.fetch(p2, m2))
			c.pc += 4
		// Input
		case 3:
			if len(c.input) == 0 {
				return -1, errNeedInput
			}
			c.set(p1, m1, c.input[0])
			c.input = c.input[1:]
			c.pc += 2
		// Output
		case 4:
			out := c.fetch(p1, m1)
			//fmt.Println("OUT:", out)
			c.pc += 2
			return out, errSuspend
		// jump-if-true
		case 5:
			if c.fetch(p1, m1) != 0 {
				c.pc = c.fetch(p2, m2)
			} else {
				c.pc += 3
			}
		// jump-if-false
		case 6:
			if c.fetch(p1, m1) == 0 {
				c.pc = c.fetch(p2, m2)
			} else {
				c.pc += 3
			}
		// less than
		case 7:
			if c.fetch(p1, m1) < c.fetch(p2, m2) {
				c.set(p3, m3, 1)
			} else {
				c.set(p3, m3, 0)
			}
			c.pc += 4
		// equals
		case 8:
			if c.fetch(p1, m1) == c.fetch(p2, m2) {
				c.set(p3, m3, 1)
			} else {
				c.set(p3, m3, 0)
			}
			c.pc += 4
		// adjust relative base
		case 9:
			c.rb += c.fetch(p1, m1)
			c.pc += 2
		// Done
		case 99:
			//fmt.Println("Done")
			return -1, nil
		default:
			panic("invalid opcode: " + strconv.Itoa(op))
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
	// part 1
	program[0] = 1

	c := newCPU(program)

	buf := make([]byte, 0)
	var done bool
	for !done {
		out, err := c.run()
		switch err {
		case nil:
			done = true

		case errSuspend:
			buf = append(buf, byte(out))

		case errNeedInput:
			fmt.Println("input")
			panic("this part shouldn't need input?")
		default:
			panic(err)
		}
	}

	var grid [][]string
	lines := strings.Split(string(buf), "\n")
	for _, l := range lines {
		fmt.Println(l)
		if len(l) > 0 {
			grid = append(grid, strings.Split(l, ""))
		}
	}

	sum := 0
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[0]); col++ {
			// anything with at least 3 neighbors is an intersection

			// TODO: need to worry about robot tile?
			if grid[row][col] != "#" {
				continue
			}

			adj := 0
			// up
			if row > 0 && grid[row-1][col] == "#" {
				adj++
			}
			// down
			if row < len(grid)-1 && grid[row+1][col] == "#" {
				adj++
			}
			// left
			if col > 0 && grid[row][col-1] == "#" {
				adj++
			}
			// right
			if col < len(grid[0])-1 && grid[row][col+1] == "#" {
				adj++
			}

			if adj == 3 {
				fmt.Printf("3? row: %d, col: %d\n", row, col)
			}
			if adj > 2 {
				sum += row * col
			}

		}
	}

	return sum
}

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
	// part 2
	program[0] = 2

	c := newCPU(program)

	// input. TODO: find much better way to express this
	mainIn := "C,A,C,B,C,A,B,B,C,A\n"
	inA := []int{'R', ',', '1', '2', ',', 'R', ',', '8', ',', 'L', ',', '8', ',', 'L', ',', '1', '2', '\n'}
	inB := []int{'L', ',', '1', '2', ',', 'L', ',', '1', '0', ',', 'L', ',', '8', '\n'}
	inC := []int{'R', ',', '8', ',', 'L', ',', '1', '0', ',', 'R', ',', '8', '\n'}
	for _, r := range mainIn {
		c.input = append(c.input, int(r))
	}
	for _, r := range inA {
		c.input = append(c.input, int(r))
	}
	for _, r := range inB {
		c.input = append(c.input, int(r))
	}
	for _, r := range inC {
		c.input = append(c.input, int(r))
	}

	// video feed yes / no
	c.input = append(c.input, []int{'n', '\n'}...)

	buf := make([]byte, 0)
	var done bool
	var lastOutput int
	for !done {
		out, err := c.run()
		switch err {
		case nil:
			return lastOutput

		case errSuspend:
			lastOutput = out

			// camera mode
			if out == '\n' && buf[len(buf)-1] == '\n' {
				fmt.Println(string(buf))
				buf = make([]byte, 0)
			}

			buf = append(buf, byte(out))

		case errNeedInput:
			fmt.Println("input")
			panic("this part shouldn't need input?")
		default:
			panic(err)
		}
	}

	return -1
}

func (c *cpu) set(i, mode, val int) {
	switch parameterMode(mode) {
	case position:
		c.state[c.state[i]] = val
	case relative:
		c.state[c.state[i]+c.rb] = val
	case immediate:
		panic("invalid mode for set")
	}
}

func (c *cpu) fetch(i int, mode int) int {

	switch parameterMode(mode) {
	case position:
		return c.state[c.state[i]]
	case relative:
		return c.state[c.state[i]+c.rb]
	case immediate:
		return c.state[i]
	default:
		panic("invalid paramater mode: " + strconv.Itoa(mode))
	}
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
	input = open("input.txt")
	fmt.Printf("B: %d\n", solveB(input))
}
