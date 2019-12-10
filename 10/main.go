package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"sort"
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
					fmt.Printf("(%d, %d)\n", y, x)
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

			if dx == 0 {
				if dy > 0 {
					slopes["down"] = true
				}
				if dy < 0 {
					slopes["up"] = true
				}
			} else if dy == 0 {
				if dx > 0 {
					slopes["right"] = true
				}
				if dx < 0 {
					slopes["left"] = true
				}
			} else {
				rat := big.NewRat(int64(dy), int64(dx))
				mod := ">"
				if dx < 0 {
					mod = "<"
				}
				slopes[rat.RatString()+mod] = true
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
	//fmt.Printf("r: %d, c: %d, %d -- %+v\n", row, col, len(slopes), slopes)
	return len(slopes)
}

// better bfs using arctan2
func bfs2(row, col int, angle float64, space [][]rune) float64 {
	//fmt.Println("Angle:", angle)

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
			//slopes[fmt.Sprintf("%f", deg)] = point
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
	fmt.Printf("LAZERING r: %d, c: %d\n", removeRow, removeCol)

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

	// all points arctan2
	// start with straight up
	// remove closest asteroid with
	// 219 too low

	prev := float64(269.9)
	for i := 0; i < 200; i++ {
		fmt.Printf("i: %d\n", i+1)
		//prev = bfs2(3, 8, prev, space)
		prev = bfs2(14, 19, prev, space)
		//fmt.Println("PREV: ", prev)
	}
	/*
		for y, r := range space {
			for x, c := range r {
				if c == '#' {
					//v := bfs2(y, x, space)
				}
			}
		}
	*/

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
