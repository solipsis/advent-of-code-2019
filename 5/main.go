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

func solveA(r io.Reader, input int) int {
	sc := bufio.NewScanner(r)

	state = make([]int, 0)

	for sc.Scan() {
		nums := strings.Split(sc.Text(), ",")
		for _, v := range nums {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			state = append(state, i)
		}
	}

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
			// hard coded input
			state[state[pc+1]] = input
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
			return result
		default:
			panic("invalid opcode: " + strconv.Itoa(op))
		}
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
	fmt.Printf("A: %d\n", solveA(input, 1))
	input = open("input.txt")
	fmt.Printf("B: %d\n", solveA(input, 5))
}
