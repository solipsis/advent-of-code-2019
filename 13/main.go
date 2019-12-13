package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type parameterMode int

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
	//fmt.Printf("input: %v\n", a.input)

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
			fmt.Println("Done")
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

	var grid [100][100]int

	c := newCPU(program)

	output := make([]int, 0)
	for {
		out, err := c.run()
		switch err {
		case nil:
			count := 0
			for _, row := range grid {
				for _, col := range row {
					if col == 2 {
						count++
					}
				}
			}
			return count
		case errSuspend:
			// new block coordinates
			output = append(output, out)
			if len(output) == 3 {
				x, y, id := output[0], output[1], output[2]
				grid[y][x] = id
				output = nil
			}
		case errNeedInput:
			// the robot is asking the color of the current tile
			//c.input = append(c.input, m[fmt.Sprintf("%d:%d", x, y)])
			panic("this part shouldn't need input?")
		default:
			panic(err)
		}
	}

	return -1
}

func printGrid(grid [25][50]int, score, loops int) {

	fmt.Println("***********************")
	fmt.Printf("loops: %d, score: %d\n", loops, score)
	for _, row := range grid {
		str := ""
		for _, col := range row {
			str += strconv.Itoa(col)
		}
		fmt.Println(str)
	}
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

	var grid [25][50]int

	c := newCPU(program)
	c.state[0] = 2 // infinite quarters

	output := make([]int, 0)
	score := 0
	x := 0
	for {

		x++

		out, err := c.run()
		switch err {
		case nil:
			return score
		case errSuspend:
			fmt.Println("output: ", x)
			// new block coordinates
			output = append(output, out)
			fmt.Println("len: ", len(output))
			if len(output) == 3 {
				x, y, id := output[0], output[1], output[2]
				output = nil // reset input buffer
				fmt.Printf("update X: %d, Y: %d, id: %d\n", x, y, id)

				// update score
				if x == -1 && y == 0 {
					score = id
					continue
				}

				// update display
				grid[y][x] = id
			}
		case errNeedInput:
			printGrid(grid, score, x)
			fmt.Println("requesting input")
			time.Sleep(time.Millisecond * 100)
			// find the x coordinate of the ball and paddle
			var bx, px int
			for _, row := range grid {
				for idx, col := range row {
					if col == 4 {
						bx = idx
					}
					if col == 3 {
						px = idx
					}
				}
			}

			if px < bx { // right
				c.input = append(c.input, 1)
			} else if px > bx { // left
				c.input = append(c.input, -1)
			} else {
				c.input = append(c.input, 0)
			}

			/*
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println(errors.New("unable to read user input: " + err.Error()))
					continue
				}
				switch input[0] {
				case 'a':
					fmt.Println("left")
					c.input = append(c.input, -1)
				case 'd':
					fmt.Println("right")
					c.input = append(c.input, 1)
				default:
					fmt.Println("neutral")
					c.input = append(c.input, 0)
				}
			*/

			//c.input = append(c.input, m[fmt.Sprintf("%d:%d", x, y)])
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
	//fmt.Printf("A: %d\n", solveA(input))
	fmt.Printf("B: %d\n", solveB(input))
}
