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

//var state []int

func (a *amplifier) init(program, input []int) {
	a.state = make([]int, len(program))
	copy(a.state, program)
	a.input = input
}

func (a *amplifier) run() (int, bool) {
	//fmt.Printf("input: %v\n", a.input)

	for {
		op := a.state[a.pc] % 100
		m1 := (a.state[a.pc] / 100) % 10
		m2 := (a.state[a.pc] / 1000) % 10
		//m3 := (a.state[a.pc] / 10000) % 10

		//fmt.Printf("op: %d, m1: %d, m2: %d, m3: %d, a.pc: %d\n", op, m1, m2, m3, pc)
		//fmt.Printf("%v\n", a.state[a.pc:pc+4])
		switch op {
		case 1:
			a.state[a.state[a.pc+3]] = a.fetch(a.pc+1, m1) + a.fetch(a.pc+2, m2)
			a.pc += 4
			break
		case 2:
			a.state[a.state[a.pc+3]] = a.fetch(a.pc+1, m1) * a.fetch(a.pc+2, m2)
			a.pc += 4
			break
		case 3:
			if len(a.input) == 0 {
				panic("No more input available")
			}
			a.state[a.state[a.pc+1]] = a.input[0]
			a.input = a.input[1:]
			a.pc += 2
		case 4:
			out := a.fetch(a.pc+1, m1)
			fmt.Println("OUT:", out)
			a.pc += 2
			return out, false
		case 5: // jump-if-true
			if a.fetch(a.pc+1, m1) != 0 {
				a.pc = a.fetch(a.pc+2, m2)
			} else {
				a.pc += 3
			}
		case 6: // jump-if-false
			if a.fetch(a.pc+1, m1) == 0 {
				a.pc = a.fetch(a.pc+2, m2)
			} else {
				a.pc += 3
			}
		case 7: // less than
			if a.fetch(a.pc+1, m1) < a.fetch(a.pc+2, m2) {
				a.state[a.state[a.pc+3]] = 1
			} else {
				a.state[a.state[a.pc+3]] = 0
			}
			a.pc += 4
		case 8: // equals
			if a.fetch(a.pc+1, m1) == a.fetch(a.pc+2, m2) {
				a.state[a.state[a.pc+3]] = 1
			} else {
				a.state[a.state[a.pc+3]] = 0
			}
			a.pc += 4

		case 99:
			fmt.Println("Done")
			return -1, true
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
	perm([]int{5, 6, 7, 8, 9}, 0, visit)

	max := 0
	for _, p := range perms {

		// re-init for each permutation set
		amplifiers := make([]*amplifier, 5)
		for x := 0; x < len(amplifiers); x++ {
			amp := &amplifier{}
			amp.init(program, []int{p[x]})
			amplifiers[x] = amp
		}
		amplifiers[0].input = append(amplifiers[0].input, 0)

		last := 0
		for i := 0; i < len(p); {
			out, done := amplifiers[i].run()

			// save prevous output of final amplifier
			if i == len(p)-1 {
				// check if this was the final cycle
				if done {
					if last > max {
						max = last
					}
					break
				}
				fmt.Println("LAST: ", out)
				last = out
			}

			if done && i == len(p)-1 {
				fmt.Println("done out", out)
				if out > max {
					max = out
					fmt.Printf("New Best: %d, perm: %v\n", max, p)
				}
				break
			}
			next := (i + 1) % len(amplifiers)
			amplifiers[next].input = append(amplifiers[next].input, out)
			i = next
		}
	}

	return max
}

type amplifier struct {
	state []int
	input []int
	pc    int
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

func (a *amplifier) fetch(i int, mode int) int {

	switch parameterMode(mode) {
	case position:
		return a.state[a.state[i]]
	case immediate:
		return a.state[i]
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
