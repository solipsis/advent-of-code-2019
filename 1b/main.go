package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(f)

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
	fmt.Println(sum)
}
