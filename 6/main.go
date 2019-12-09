package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type node struct {
	id     string
	orbits map[string]*node
}

var galaxy = make(map[string]*node)

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	galaxy["COM"] = &node{id: "COM", orbits: nil}
	for sc.Scan() {
		arr := strings.Split(sc.Text(), ")")
		a, b := arr[0], arr[1]

		if galaxy[a] == nil {
			galaxy[a] = &node{id: a, orbits: make(map[string]*node)}
		}
		if galaxy[b] == nil {
			galaxy[b] = &node{id: b, orbits: make(map[string]*node)}
		}

		galaxy[b].orbits[a] = galaxy[a]
	}

	total := 0
	for _, n := range galaxy {
		total += dfs(n)
	}

	return total
}

func dfs(n *node) int {
	if n == nil {
		return 0
	}
	sum := 0

	for _, sub := range n.orbits {
		sum += 1 + dfs(sub)
	}
	return sum
}

func main() {
	input := open("input.txt")
	fmt.Printf("A: %d\n", solveA(input))
	//input = open("input.txt")
	//fmt.Printf("B: %d\n", solveB(input))
}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}
