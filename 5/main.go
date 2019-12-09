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

func solveA(r io.Reader) int {
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

	pc := 0
	for {
		//fmt.Println("----------------------------------------")
		//fmt.Printf("%v\n", state)
		//t.Println(pc)

		op := state[pc] % 100
		m1 := (state[pc] / 100) % 10
		m2 := (state[pc] / 1000) % 10
		m3 := (state[pc] / 10000) % 10

		fmt.Printf("op: %d, m1: %d, m2: %d, m3: %d, pc: %d\n", op, m1, m2, m3, pc)
		fmt.Printf("%v\n", state[pc:pc+4])
		switch op {
		case 1:
			//state[state[pc+3]] = state[state[pc+1]] + state[state[pc+2]]
			state[state[pc+3]] = fetch(pc+1, m1) + fetch(pc+2, m2)
			//fmt.Printf("add A: %d, B: %d, C: %d\n", fetch(pc+1, m1), fetch(pc+2, m2), state[pc+3])
			pc += 4
			break
		case 2:
			//state[state[pc+3]] = state[state[pc+1]] * state[state[pc+2]]
			state[state[pc+3]] = fetch(pc+1, m1) * fetch(pc+2, m2)
			pc += 4
			break
		case 3:
			// hard coded input
			state[state[pc+1]] = 1
			pc += 2
		case 4:
			fmt.Println("OUT:", fetch(pc+1, m1))
			pc += 2
		case 99:
			return state[0]
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

/*
func solveB(r io.Reader) int {
	sc := bufio.NewScanner(r)
	initial := make([]int, 0)

	for sc.Scan() {
		nums := strings.Split(sc.Text(), ",")
		for _, v := range nums {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			initial = append(initial, i)
		}
	}

	for noun := 1; noun < len(initial); noun++ {
		for verb := 1; verb < len(initial); verb++ {
			state := make([]int, len(initial))
			copy(state, initial)
			if execProgram(state, noun, verb) == 19690720 {
				return 100*noun + verb
			}
		}
	}
	return -1
}

func execProgram(state []int, noun, verb int) int {
	state[1] = noun
	state[2] = verb
	for pc := 0; ; pc += 4 {
		switch state[pc] {
		case 1:
			state[state[pc+3]] = state[state[pc+1]] + state[state[pc+2]]
			break
		case 2:
			state[state[pc+3]] = state[state[pc+1]] * state[state[pc+2]]
			break
		case 99:
			return state[0]
		default:
			panic("invalid opcode")
		}
	}
}

*/
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
