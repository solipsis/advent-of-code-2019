package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type node struct {
	id      string
	orbits  map[string]*node
	links   map[string]*node
	visited bool
}

var galaxy map[string]*node

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	galaxy = make(map[string]*node)
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

func solveB(r io.Reader) int {
	sc := bufio.NewScanner(r)

	galaxy = make(map[string]*node)
	galaxy["COM"] = &node{id: "COM", links: make(map[string]*node)}
	for sc.Scan() {
		arr := strings.Split(sc.Text(), ")")
		a, b := arr[0], arr[1]

		if galaxy[a] == nil {
			galaxy[a] = &node{id: a, links: make(map[string]*node)}
		}
		if galaxy[b] == nil {
			galaxy[b] = &node{id: b, links: make(map[string]*node)}
		}

		// link in both directions
		galaxy[b].links[a] = galaxy[a]
		galaxy[a].links[b] = galaxy[b]
	}

	return bfs(galaxy["YOU"], "SAN") - 2
}

func bfs(n *node, target string) int {

	queue := []*node{n, nil}
	depth := 0
	for len(queue) > 1 {
		// pop top item
		next := queue[0]
		queue = queue[1:]

		// check if depth has increased
		if next == nil {
			depth++
			// depth marker
			queue = append(queue, nil)
			continue
		}

		fmt.Printf("visiting: %v, depth: %d\n", (*next).id, depth)

		if next.id == target {
			return depth
		}

		next.visited = true // prevent loops

		// visit children
		for _, l := range next.links {
			if !l.visited {
				queue = append(queue, l)
			}
		}
	}

	return -1
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
