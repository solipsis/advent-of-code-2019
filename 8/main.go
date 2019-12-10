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
	return result
}

type color int

const (
	black color = iota
	white
	transparent
)

func solveB(r io.Reader) int {

	sc := bufio.NewScanner(r)

	width := 25
	height := 6
	image := make([]layer, 0)

	var row, col int
	for sc.Scan() {
		arr := strings.Split(sc.Text(), "")
		layer := newLayer(width, height)
		for _, v := range arr {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			layer[row][col] = i
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
			}
		}
	}

	complete := newLayer(width, height)

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			for lay := 0; lay < len(image); lay++ {
				if color(image[lay][row][col]) == black || color(image[lay][row][col]) == white {
					complete[row][col] = image[lay][row][col]
					break
				}
			}
		}
	}

	for _, r := range complete {
		for _, c := range r {
			fmt.Printf("%d", c)
		}
		fmt.Println()
	}

	return -1
}

func main() {
	//input := open("input.txt")
	//fmt.Printf("A: %d\n", solveA(input))

	input := open("input.txt")
	fmt.Printf("B: %d\n", solveB(input))
}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}
