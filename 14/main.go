package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type recipe struct {
	outID  string
	outAmt int
	inputs map[string]int
}

type state struct {
	toConsume int
	nextState map[string]int
}

var states map[string]*state

var supply map[string]int
var recipes map[string]recipe
var oreConsumed int
var fuelProduced int

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	supply = make(map[string]int)
	recipes = make(map[string]recipe)
	states = make(map[string]*state)
	oreConsumed = 0

	for sc.Scan() {
		lr := strings.Split(sc.Text(), " => ")
		left, right := lr[0], lr[1]

		// output
		pair := strings.Split(right, " ")
		amt, err := strconv.Atoi(pair[0])
		if err != nil {
			panic(err)
		}

		rec := recipe{
			outID:  pair[1],
			outAmt: amt,
			inputs: make(map[string]int),
		}

		// inputs
		inputs := strings.Split(left, ", ")
		for _, in := range inputs {
			arr := strings.Split(in, " ")
			//fmt.Println("in arrr", arr)
			count, err := strconv.Atoi(arr[0])
			if err != nil {
				panic(err)
			}
			inName := arr[1]

			rec.inputs[inName] = count
		}
		recipes[pair[1]] = rec
	}

	fmt.Printf("%+v\n", recipes)

	x := 0
	for oreConsumed < 1000000000000 {

		x++
		if x%10000 == 0 {
			fmt.Println(x, oreConsumed, fmt.Sprintf("%02f", (float64(oreConsumed)/1000000000000)*100))
		}

		//prevOre := oreConsumed

		/*
			// check if we have already seen this
			stateStr := stateToString()
			if states[stateStr] != nil {
				//fmt.Println("old")
				supply = states[stateStr].nextState
				oreConsumed += states[stateStr].toConsume
				continue
			}

			//fmt.Println("new", len(states))
			//fmt.Println("new", stateStr, len(states))
			if x%1000 == 0 {
				fmt.Println("new", x, len(states))
			}
		*/
		// new state
		produce("FUEL", 1)
		//newSupply := copySupply()
		//s := &state{
		//toConsume: oreConsumed - prevOre,
		//nextState: newSupply,
		//}
		//states[stateStr] = s

	}

	/*
		for x := 1000; x <= 100000; x += 10000 {
			oreConsumed = 0
			produce("FUEL", x)
			//fmt.Printf("Consumed: %d, diff: %d\n", oreConsumed, prev-oreConsumed)
			//prev = oreConsumed
			fmt.Printf("Fuel: %d, ore per fuel: %d\n", x, oreConsumed/x)
		}
	*/

	/*
		produce("FUEL", 1)
		fmt.Printf("Consumed: %d\n", oreConsumed)
		prev := oreConsumed
		for x := 0; x < 50; x++ {
			oreConsumed = 0
			produce("FUEL", 1)
			fmt.Printf("Consumed: %d, diff: %d\n", oreConsumed, prev-oreConsumed)
			prev = oreConsumed
		}
	*/

	//}

	fmt.Println("fuelProduced: ", x-1, "consumed:", oreConsumed)
	return x - 1
}

func copySupply() map[string]int {
	m := make(map[string]int)
	for k, v := range supply {
		m[k] = v
	}
	return m
}

func stateToString() string {
	/*
		sorted := make([]string, len(supply))
		x := 0
		for k := range supply {
			sorted[x] = k
			x++
		}

		sort.Slice(sorted, func(x, y int) bool {
			return sorted[x] < sorted[y]
		})
	*/
	//return fmt.Sprintf("%v", sorted)
	return fmt.Sprintf("%v", supply)
}

// too low 4967117

func produce(outID string, outAmt int) {
	//fmt.Printf("Producing: %d %s %v\n", outAmt, outID, recipes[outID])
	//fmt.Printf("Current Supply: %v\n", supply)
	/*
		if supply[outID] >= outAmt {
			panic("no need to produce why am I here?")
			return
		}
	*/

	produced := 0

	rec := recipes[outID]
	for produced < outAmt {

		for inID, inAmt := range rec.inputs {
			//fmt.Printf("Input: %d %s\n", inAmt, inID)

			// special ORE case
			if inID == "ORE" {
				//supply["ORE"] += inAmt
				//fmt.Println("Producing ORE")
				supply[outID] += rec.outAmt
				oreConsumed += inAmt
				return
			}

			if supply[inID] >= inAmt {
				supply[inID] -= inAmt
				continue
			}
			for supply[inID] < inAmt {
				//fmt.Println("not enough producing more")
				produce(inID, inAmt-supply[inID])
			}
			supply[inID] -= inAmt
		}

		if outID == "FUEL" {
			return
		}

		supply[outID] += rec.outAmt
		produced += rec.outAmt
	}
}

func solveB(r io.Reader) int {
	return -1
}

func main() {
	input := open("input.txt")
	fmt.Printf("A: %d\n", solveA(input))
	//fmt.Printf("B: %d\n", solveB(input))
}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}
