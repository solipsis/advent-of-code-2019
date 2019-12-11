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

	type dir struct {
		x, y int
	}
	// left up right down
	dirs := []dir{
		dir{-1, 0},
		dir{0, -1},
		dir{1, 0},
		dir{0, 1},
	}

	//grid := make[][]j
	m := make(map[string]int)
	var x, y int

	c := newCPU(program)

	//c.input = append(c.input, 2)
	paint := true
	curDir := 1
	for {
		out, err := c.run()
		switch err {
		case nil:
			return len(m)
		case errSuspend:
			// the robot has output a color and needs to move
			fmt.Println("Out:", out)
			if paint {
				paint = false
				if out == 0 {
					m[fmt.Sprintf("%d:%d", x, y)] = 0
				} else if out == 1 {
					m[fmt.Sprintf("%d:%d", x, y)] = 1
				}
			} else {
				paint = true
				// move
				if out == 0 {
					curDir = (curDir - 1 + len(dirs)) % len(dirs)
				} else {
					curDir = (curDir + 1) % len(dirs)
				}
				x += dirs[curDir].x
				y += dirs[curDir].y
			}

		case errNeedInput:
			// the robot is asking the color of the current tile
			c.input = append(c.input, m[fmt.Sprintf("%d:%d", x, y)])
		default:
			panic(err)
		}
	}

	return -1
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

	type dir struct {
		x, y int
	}
	// left up right down
	dirs := []dir{
		dir{-1, 0},
		dir{0, -1},
		dir{1, 0},
		dir{0, 1},
	}

	m := make(map[string]int)
	var x, y int

	c := newCPU(program)

	//c.input = append(c.input, 2)
	paint := true
	curDir := 1
	for {
		out, err := c.run()
		switch err {
		case nil:
			return len(m)
		case errSuspend:
			// the robot has output a color and needs to move
			fmt.Println("Out:", out)
			if paint {
				paint = false
				if out == 0 {
					m[fmt.Sprintf("%d:%d", x, y)] = 0
				} else if out == 1 {
					m[fmt.Sprintf("%d:%d", x, y)] = 1
				}
			} else {
				paint = true
				// move
				if out == 0 {
					curDir = (curDir - 1 + len(dirs)) % len(dirs)
				} else {
					curDir = (curDir + 1) % len(dirs)
				}
				x += dirs[curDir].x
				y += dirs[curDir].y
			}

		case errNeedInput:
			// the robot is asking the color of the current tile
			c.input = append(c.input, m[fmt.Sprintf("%d:%d", x, y)])
		default:
			panic(err)
		}
	}

	return -1
}

//func pointStr(x, y) string {
//return
//}

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
}
