package main

import (
	"bufio"
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
)

var mode parameterMode = position

var state []int

func execProgram(program, input []int) int {
	fmt.Printf("Program: %v\n", program)
	fmt.Printf("Input: %v\n", input)

	state = make([]int, len(program))
	copy(state, program)

	result := 0
	pc := 0
	for {
		op := state[pc] % 100
		m1 := (state[pc] / 100) % 10
		m2 := (state[pc] / 1000) % 10
		//m3 := (state[pc] / 10000) % 10

		//fmt.Printf("op: %d, m1: %d, m2: %d, m3: %d, pc: %d\n", op, m1, m2, m3, pc)
		//fmt.Printf("%v\n", state[pc:pc+4])
		switch op {
		case 1:
			state[state[pc+3]] = fetch(pc+1, m1) + fetch(pc+2, m2)
			pc += 4
			break
		case 2:
			state[state[pc+3]] = fetch(pc+1, m1) * fetch(pc+2, m2)
			pc += 4
			break
		case 3:
			if len(input) == 0 {
				panic("No more input available")
			}
			state[state[pc+1]] = input[0]
			input = input[1:]
			pc += 2
		case 4:
			fmt.Println("OUT:", fetch(pc+1, m1))
			result = fetch(pc+1, m1)
			pc += 2
		case 5: // jump-if-true
			if fetch(pc+1, m1) != 0 {
				pc = fetch(pc+2, m2)
			} else {
				pc += 3
			}
		case 6: // jump-if-false
			if fetch(pc+1, m1) == 0 {
				pc = fetch(pc+2, m2)
			} else {
				pc += 3
			}
		case 7: // less than
			if fetch(pc+1, m1) < fetch(pc+2, m2) {
				state[state[pc+3]] = 1
			} else {
				state[state[pc+3]] = 0
			}
			pc += 4
		case 8: // equals
			if fetch(pc+1, m1) == fetch(pc+2, m2) {
				state[state[pc+3]] = 1
			} else {
				state[state[pc+3]] = 0
			}
			pc += 4

		case 99:
			fmt.Println("Done")
			return result
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

	perms := [][]int{}
	visit := func(perm []int) {
		cpy := make([]int, len(perm))
		copy(cpy, perm)
		perms = append(perms, cpy)
	}
	perm([]int{0, 1, 2, 3, 4}, 0, visit)

	max := 0
	out := 0
	for _, p := range perms {
		fmt.Printf("perm: %v\n", p)
		// first execution
		out = execProgram(program, []int{p[0], 0})
		for _, v := range p[1:] {
			out = execProgram(program, []int{v, out})
		}
		if out > max {
			max = out
		}
	}

	return max
}

func perm(arr []int, i int, visit func(perm []int)) {
	if i >= len(arr) {
		visit(arr)
		return
	}
	for j := i; j < len(arr); j++ {
		arr[i], arr[j] = arr[j], arr[i]
		perm(arr, i+1, visit)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func fetch(i int, mode int) int {

	switch parameterMode(mode) {
	case position:
		return state[state[i]]
	case immediate:
		return state[i]
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
