package main

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
)

func solveA(r io.Reader) int {

	// BFS from each Asteroid and store slope?
	// any asteroids at same slope visited later don't count?

	sc := bufio.NewScanner(r)

	space := make([][]rune, 0)

	for sc.Scan() {
		arr := strings.Split(sc.Text(), "")

		row := make([]rune, 0)
		for _, v := range arr {
			row = append(row, rune(v[0]))
		}
		space = append(space, row)
	}

	best := 0
	for y, r := range space {
		for x, c := range r {
			if c == '#' {
				v := bfs(y, x, space)
				if v > best {
					best = v
				}
			}
		}
	}

	// for each xy
	// append 8 surrounding

	return best
}

func bfs(row, col int, space [][]rune) int {

	type point struct {
		r, c int
	}

	visited := make(map[point]bool)
	slopes := make(map[string]bool)

	queue := []point{point{row, col}}
	for len(queue) > 0 {
		// pop
		next := queue[0]
		queue = queue[1:]

		r, c := next.r, next.c

		if r < 0 || r >= len(space) || c < 0 || c >= len(space[0]) {
			continue
		}
		if visited[next] {
			continue
		}
		visited[next] = true

		// calculate slope if this is an asteroid
		if space[r][c] == '#' {
			dy := r - row
			dx := c - col

			if dx != 0 {
				rat := big.NewRat(int64(dy), int64(dx))

				if dy == 0 {
					if dx > 0 {
						slopes["right"] = true
					} else if dx < 0 {
						slopes["left"] = true
					}
				} else {
					slopes[rat.RatString()] = true
				}
			} else {
				if dy > 0 {
					slopes["down"] = true
				} else if dy < 0 {
					slopes["up"] = true
				}
			}
		}

		// search 8 surrounding squares
		points := []point{
			point{r - 1, c - 1},
			point{r - 1, c},
			point{r - 1, c + 1},
			point{r, c - 1},
			point{r, c + 1},
			point{r + 1, c - 1},
			point{r + 1, c},
			point{r + 1, c + 1},
		}
		queue = append(queue, points...)
	}
	fmt.Printf("r: %d, c: %d, %d -- %+v\n", row, col, len(slopes), slopes)
	return len(slopes)
}

func main() {
	input := open("input.txt")
	fmt.Printf("A: %d\n", solveA(input))

}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}
