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

var supply map[string]int
var recipes map[string]recipe
var oreConsumed int

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	supply = make(map[string]int)
	recipes = make(map[string]recipe)
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
	//for oreConsumed < 1000000000000 {
	//produce("FUEL", 1)
	//}

	return oreConsumed
}

// too low 4967117

func produce(outID string, outAmt int) {
	fmt.Printf("Producing: %d %s %v\n", outAmt, outID, recipes[outID])
	fmt.Printf("Current Supply: %v\n", supply)
	/*
		if supply[outID] >= outAmt {
			panic("no need to produce why am I here?")
			return
		}
	*/

	rec := recipes[outID]

	for inID, inAmt := range rec.inputs {
		fmt.Printf("Input: %d %s\n", inAmt, inID)

		// special ORE case
		if inID == "ORE" {
			//supply["ORE"] += inAmt
			fmt.Println("Producing ORE")
			supply[outID] += rec.outAmt
			oreConsumed += inAmt
			return
		}

		if supply[inID] >= inAmt {
			supply[inID] -= inAmt
			continue
		}
		for supply[inID] < inAmt {
			fmt.Println("not enough producing more")
			produce(inID, 1)
		}
		supply[inID] -= inAmt
	}
	supply[outID] += rec.outAmt
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
