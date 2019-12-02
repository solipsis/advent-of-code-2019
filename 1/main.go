package main

import (
	"bufio"
	"io"
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

func main() {}
