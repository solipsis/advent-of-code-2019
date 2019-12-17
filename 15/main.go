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
			fmt.Println("Done")
			return -1, nil
		default:
			panic("invalid opcode: " + strconv.Itoa(op))
		}
	}

}

var grid [][]int

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

	size := 50
	grid = make([][]int, size)
	for i := range grid {
		grid[i] = make([]int, size)
	}
	//x, y := 5, 5

	c := newCPU(program)

	/*
		c.input = append(c.input, 4)
		out, err := c.run()
		fmt.Println(err)
		fmt.Println(out)

		c.input = append(c.input, 3)
		out, err = c.run()
		fmt.Println(err)
		fmt.Println(out)
	*/
	dfs(c, 25, 25, -1)
	grid[25][25] = 8
	grid[43][41] = 2
	printGrid(grid)

	// part 1
	depth := bfs(grid, 25, 25)

	// part 2
	depth = bfs(grid, 43, 41)

	return depth
}

var dirXY = map[int][]int{
	1: {0, -1},
	2: {0, 1},
	3: {-1, 0},
	4: {1, 0},
}

var visited = make(map[string]bool)

func dfs(c *cpu, x, y, prevDir int) {

	//vstr :=
	//fmt.Println("***********************************************")
	//fmt.Printf("x: %d, y: %d, prev: %d\n", x, y, prevDir)
	if visited[fmt.Sprintf("(x:%d,y:%d)", x, y)] {
		panic("visiting same node twice, why?")
	}
	visited[fmt.Sprintf("(x:%d,y:%d)", x, y)] = true
	grid[y][x] = 1
	//printGrid(grid)

	// NORTH SOUTH WEST EAST
	dirs := []int{1, 2, 3, 4}

	for _, dir := range dirs {

		// initiate the new movement
		//fmt.Printf("Input %d\n", dir)

		// skip any square we have already visited
		nx := x + dirXY[dir][0]
		ny := y + dirXY[dir][1]
		if grid[ny][nx] != 0 {
			continue
		}
		c.input = append(c.input, dir)

		out, err := c.run()
		//fmt.Printf("Response %d\n", out)
		switch err {
		case nil:

		case errSuspend:
			// movement response
			switch out {
			// hit a wall
			case 0:
				//fmt.Printf("marking wall: wx: %d, wy: %d, dir: %d\n", nx, ny, dir)
				grid[ny][nx] = 9
			// successful move
			case 1:
				if grid[ny][nx] == 0 {
					//fmt.Printf("recursing to nx: %d, ny, dir: %d: %d\n", nx, ny, dir)
					dfs(c, nx, ny, dir)

					//fmt.Printf("undoing move to nx: %d, ny: %d\n", nx, ny)
					//c.run()
				}
			case 2:
				fmt.Printf("OXYGEN: x:%d, y:%d\n", nx, ny)
				grid[ny][nx] = 2
				dfs(c, nx, ny, dir)
			default:
				panic("blah")

			}

		case errNeedInput:
			// the robot is asking the color of the current tile
			//c.input = append(c.input, m[fmt.Sprintf("%d:%d", x, y)])
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
		//if

	}
	return depth
}

func printGrid(g [][]int) {
	//fmt.Printf("%v\n", visited)
	for _, r := range grid {
		for _, c := range r {
			fmt.Printf("%d", c)
		}
		fmt.Println()
	}

}

func solveB(r io.Reader) int {

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
	//fmt.Printf("B: %d\n", solveB(input))
}
