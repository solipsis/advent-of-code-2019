package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

// better bfs using arctan2
func bfs2(row, col int, angle float64, space [][]rune) float64 {

	// BFS to find the first node we encounter with each unique angle
	// relative to the starting position
	// Then find the closest node with an angle greater than
	// the given angle

	type point struct {
		r, c int
	}

	visited := make(map[point]bool)
	slopes := make(map[float64]point)

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

		if space[r][c] == '#' {
			dy := r - row
			dx := c - col

			// convert to positive degrees
			rad := math.Atan2(float64(dy), float64(dx))
			deg := rad * (180 / math.Pi)
			if deg < 0 {
				deg += 360
			}

			// only add the closest node we find with this angle
			if _, ok := slopes[deg]; !ok {
				slopes[deg] = next
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

	sorted := make([]float64, 0)
	for k := range slopes {
		sorted = append(sorted, k)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	next := 0
	for idx, v := range sorted {
		if v > angle {
			next = idx
			break
		}
	}
	removeRow := slopes[sorted[next]].r
	removeCol := slopes[sorted[next]].c

	space[removeRow][removeCol] = '.'
	fmt.Printf("row: %d, col: %d\n", removeRow, removeCol)

	//fmt.Printf("%v\n", sorted)
	//fmt.Printf("r: %d, c: %d, %d -- %+v\n", row, col, len(slopes), slopes)
	return sorted[next]
}

func solveB(r io.Reader) int {
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

	// angle right before straight up in this coordinate space
	prev := float64(269.9)

	// previous puzzle output
	space[14][19] = 'X'

	for i := 0; i < 200; i++ {
		prev = bfs2(14, 19, prev, space)
	}

	return -1
}

func main() {
	input := open("input.txt")
	//fmt.Printf("A: %d\n", solveA(input))
	//input := open("input.txt")
	fmt.Printf("B: %d\n", solveB(input))

}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}
