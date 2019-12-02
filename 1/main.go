package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	sum := 0
	for sc.Scan() {
		i, err := strconv.Atoi(sc.Text())
		if err != nil {
			panic(err)
		}

		sum += i/3 - 2
	}
	return sum
}

func solveB(r io.Reader) int {
	sc := bufio.NewScanner(r)

	sum := 0
	for sc.Scan() {
		i, err := strconv.Atoi(sc.Text())
		if err != nil {
			panic(err)
		}

		v := i
		for {
			v = (v / 3) - 2
			if v < 0 {
				break
			}
			sum += v
		}
	}
	return sum
}

func main() {
	input := open("input.txt")
	fmt.Printf("A: %d\n", solveA(input))
	input = open("input.txt")
	fmt.Printf("B: %d\n", solveB(input))
}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}
