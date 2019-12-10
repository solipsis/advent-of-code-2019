package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type layer [][]int

func newLayer(w, h int) layer {
	rows := make([][]int, h)
	for i := range rows {
		rows[i] = make([]int, w)
	}
	return layer(rows)
}

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	width := 25
	height := 6
	image := make([]layer, 0)
	result := -1

	var row, col int
	for sc.Scan() {
		arr := strings.Split(sc.Text(), "")
		//fmt.Println(len(arr))
		//fmt.Println(25 * 6)
		layer := newLayer(width, height)
		var c0, c1, c2 int
		fewest := 9999999
		for _, v := range arr {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			layer[row][col] = i

			if i == 0 {
				c0++
			}
			if i == 1 {
				c1++
			}
			if i == 2 {
				c2++
			}

			col++
			if col == width {
				col = 0
				row++
			}
			if row == height {
				// prepare next layer
				row = 0
				image = append(image, layer)
				layer = newLayer(width, height)

				if c0 < fewest {
					fewest = c0
					result = c1 * c2
				}
				// reset counters
				c0, c1, c2 = 0, 0, 0
			}
			fmt.Printf("row: %d, col: %d\n", row, col)
		}
	}

	/*
		best := 99999999
		bestIdx := -1
		for i, layer := range image {

			count := 0
			for _, r := range image {
				for _, c := range r {
					if c == 0 {
						count++
					}
				}
				if count < best {
					best = count
					bestIdx = i
				}
			}
		}

		// # of 1's times # of 2's
		var c1, c2 int
		for _, c := range image[bestIdx] {
			if c == 1 {
				c1++
			}
			if c == 2 {
				c2++
			}
		}
		return c1 * c2
	*/
	return result
}

/*
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
*/

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
