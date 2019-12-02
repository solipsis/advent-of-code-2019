package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	state := make([]int, 0)

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
