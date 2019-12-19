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

	c := newCPU(program)

	buf := make([]byte, 0)
	var done bool
	for !done {
		out, err := c.run()
		switch err {
		case nil:
			done = true

		case errSuspend:
			//fmt.Println("suspend, out:", out)
			buf = append(buf, byte(out))
			//fmt.Printf("%s", string(rune(out)))
			/*
				// movement response
				switch out {
				// hit a wall
				case 0:
				// successful move
				case 1:
				case 2:
				default:
					panic("blah")
				}
			*/

		case errNeedInput:
			fmt.Println("input")
			panic("this part shouldn't need input?")
		default:
			panic(err)
		}
	}

	// 6476 too high
	var grid [][]string
	lines := strings.Split(string(buf), "\n")
	for _, l := range lines {
		fmt.Println(l)
		if len(l) > 0 {
			grid = append(grid, strings.Split(l, ""))
		}
	}

	// input
	mainIn := strings.Split("C,A,C,B,C,A,B,B,C,A\n", "")
	inA := strings.Split("R12,R8,L8,L12\n", "")
	inB := strings.Split("L12,L10,L8\n", "")
	inC := strings.Split("R8,L10,R8\n", "")

	//fmt.Printf("%v\n", grid)
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

	fmt.Println("***********************************8")
	//fmt.Printf("%v\n", grid)

	return sum
}

func solveB(r io.Reader) int {
	return -1
}

var dirXY = map[int][]int{
	1: {0, -1},
	2: {0, 1},
	3: {-1, 0},
	4: {1, 0},
}

var visited map[string]bool

func dfs(c *cpu, x, y, prevDir int) {

	if visited[fmt.Sprintf("(x:%d,y:%d)", x, y)] {
		panic("visiting same node twice, why?")
	}
	visited[fmt.Sprintf("(x:%d,y:%d)", x, y)] = true
	grid[y][x] = 1

	// NORTH SOUTH WEST EAST
	dirs := []int{1, 2, 3, 4}

	for _, dir := range dirs {

		// skip any square we have already visited
		nx := x + dirXY[dir][0]
		ny := y + dirXY[dir][1]
		if grid[ny][nx] != 0 {
			continue
		}
		c.input = append(c.input, dir)

		out, err := c.run()
		switch err {
		case nil:

		case errSuspend:
			// movement response
			switch out {
			// hit a wall
			case 0:
				grid[ny][nx] = 9
			// successful move
			case 1:
				if grid[ny][nx] == 0 {
					dfs(c, nx, ny, dir)
				}
			case 2:
				fmt.Printf("OXYGEN: x:%d, y:%d\n", nx, ny)
				grid[ny][nx] = 2
				dfs(c, nx, ny, dir)
			default:
				panic("blah")
			}

		case errNeedInput:
			panic("this part shouldn't need input?")
		default:
			panic(err)
		}
	}

	if prevDir == -1 {
		return
	}
	// Undo the move that got us here
	var reverse int
	switch prevDir {
	case 1:
		reverse = 2
	case 2:
		reverse = 1
	case 3:
		reverse = 4
	case 4:
		reverse = 3
	default:
		panic("how did this happen")
	}
	//fmt.Println("REVERSE")
	c.input = append(c.input, reverse)
	c.run()
}

func bfs(g [][]int, row, col int) int {

	type point struct {
		row, col int
	}

	visited := make(map[string]bool)

	depth := 0

	// initial point + depth marker
	queue := []*point{&point{row, col}, nil}
	for len(queue) > 1 {
		// pop
		next := queue[0]
		queue = queue[1:]

		// We have hit a depth marker
		if next == nil {
			depth++
			if queue[0] == nil {
				return depth
			}
			queue = append(queue, nil)
			continue
		}

		// we found the end
		if grid[next.row][next.col] == 2 {
			//return depth
		}

		dirs := [][]int{
			[]int{0, -1},
			[]int{0, 1},
			[]int{-1, 0},
			[]int{1, 0},
		}
		// enqueue children
		for _, dir := range dirs {
			nr := next.row + dir[0]
			nc := next.col + dir[1]

			pair := fmt.Sprintf("row%s:col%s", nr, nc)
			if visited[pair] || grid[nr][nc] == 9 || grid[nr][nc] == 0 {
				continue
			}
			queue = append(queue, &point{nr, nc})
			visited[pair] = true
		}

	}
	return depth
}

func printGrid(g [][]int) {
	//fmt.Printf("%v\n", visited)
	for _, r := range grid {
		for _, c := range r {
			fmt.Printf("%s", string(rune(c)))
		}
		fmt.Println()
	}

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
