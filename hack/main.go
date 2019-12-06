package main

import (
	"fmt"
	"strconv"
	"strings"
)

var base [9][9]int
var preset = make(map[int]bool)
var debug int

//var priority
type narrow struct {
	row  int
	col  int
	vals []int
}

func main() {

	// apply orignal constraints
	base[0][8] = 1
	base[1][1] = 1
	base[1][2] = 2
	base[2][6] = 2
	base[3][8] = 2
	base[4][1] = 2
	base[6][6] = 1
	base[6][7] = 2
	base[7][0] = 1
	base[7][5] = 2
	base[8][0] = 2
	base[8][3] = 1
	preset[8] = true
	preset[11] = true
	preset[12] = true
	preset[26] = true
	preset[38] = true
	preset[41] = true
	preset[66] = true
	preset[67] = true
	preset[70] = true
	preset[75] = true
	preset[80] = true
	preset[83] = true

	// subranges for values I was able to figure out with pen and paper
	priority := make([]narrow, 8)
	priority[0] = narrow{row: 0, col: 4, vals: []int{2, 3, 4}}
	priority[1] = narrow{row: 2, col: 1, vals: []int{7, 8, 9}}
	priority[2] = narrow{row: 4, col: 8, vals: []int{3, 4, 5}}
	priority[3] = narrow{row: 5, col: 1, vals: []int{3, 4, 5}}
	priority[4] = narrow{row: 7, col: 3, vals: []int{3, 4, 5}}
	priority[5] = narrow{row: 7, col: 7, vals: []int{6, 7, 8, 9}}
	priority[6] = narrow{row: 8, col: 6, vals: []int{6, 7, 8, 9}}
	priority[7] = narrow{row: 0, col: 3, vals: []int{3, 4, 5, 6, 7}}

	// uncomment to generate giant contraint blocks
	/*
		// programattically generate constraints so I don't make a type
		equations := []string{
			"B9 + B8 + C1 + H4 + H4 = 23",
			"A5 + D7 + I5 + G8 + B3 + A5 = 19",
			"I2 + I3 + F2 + E9 = 15",
			"I7 + H8 + C2 + D9 = 26",
			"I6 + A5 + I3 + B8 + C3 = 20",
			"I7 + D9 + B6 + A8 + A3 + C4 = 27",
			"C7 + H9 + I7 + B2 + H8 + G3 = 31",
			"D3 + I8 + A4 + I6 = 27",
			"F5 + B8 + F8 + I7 + F1 = 33",
			"A2 + A8 + D7 + E4 = 21",
			"C1 + I4 + C2 + I1 + A4 = 20",
			"F8 + C1 + F6 + D3 + B6 = 25",
		}
		for _, eq := range equations {
			parseEq(eq)
		}
	*/

	results := make(chan string) // intermediate debug values
	done := make(chan string)

	// start each thread with a different first box
	for i := 3; i <= 9; i++ {
		gridCpy := make([][]int, 9)
		base[0][0] = i
		for j := 0; j < 9; j++ {
			gridCpy[j] = make([]int, 9)
			copy(gridCpy[j], base[j][:])
		}
		go solve(0, 1, priority, results, done, gridCpy)
	}

	// aggregate debug results on a single thread and wait for a solution
	for {
		select {
		case res := <-results:
			fmt.Println("\n\n----------------------------------------")
			fmt.Println(res)
		case answer := <-done:
			fmt.Println("Done")
			fmt.Println("----------------------------------------")
			fmt.Println(answer)
			// too lazy to clean up properly
			close(done)
			close(results)
			panic(answer)
		}
	}

	fmt.Println("FAIL")
}

func solve(r, c int, priority []narrow, results, done chan string, grid [][]int) {

	// print intermediate solutions for debugging
	debug++
	if debug%10000000 == 0 {
		select {
		case <-done:
			return
		case results <- printSol(grid):
		}
	}

	// recursive base case
	if r == 9 {
		if equationValidate(grid) {
			done <- printSol(grid)
		}
		return
	}

	// make sure all values that are bounded are set early as possible
	// if we haven't exhausted the ranges we know about, recurse on those first
	if len(priority) > 0 {
		p := priority[0]
		for _, v := range p.vals {
			if rowOK(p.row, v, grid) && colOK(p.col, v, grid) && boxOK(p.row, p.col, v, grid) {
				grid[p.row][p.col] = v
				if equationOK(grid) {
					solve(r, c, priority[1:], results, done, grid)
				}
			}
			grid[p.row][p.col] = 0
		}
		return
	}

	// is this square preset? if so skip
	if preset[r*10+c] {
		if c == 8 {
			solve(r+1, 0, priority, results, done, grid)
		} else {
			solve(r, c+1, priority, results, done, grid)
		}
		return
	}

	// if this is a priority square skip cause already set
	if grid[r][c] != 0 {
		if c == 8 {
			solve(r+1, 0, priority, results, done, grid)
		} else {
			solve(r, c+1, priority, results, done, grid)
		}
		return
	}

	// place values 1-9 in this square
	// go in reverse order to eagerly fail on on equation pruning bounds
	for v := 9; v >= 1; v-- {
		if rowOK(r, v, grid) && colOK(c, v, grid) && boxOK(r, c, v, grid) {
			grid[r][c] = v
			if !equationOK(grid) {
				continue
			}

			// this is valid, move to the next column, or row if this is the last column
			if c == 8 {
				solve(r+1, 0, priority, results, done, grid)
			} else {
				solve(r, c+1, priority, results, done, grid)
			}
		}
	}
	grid[r][c] = 0

}

func rowOK(row, val int, grid [][]int) bool {
	for col := 0; col < 9; col++ {
		if grid[row][col] == val {
			return false
		}
	}
	return true
}

func colOK(col, val int, grid [][]int) bool {
	for row := 0; row < 9; row++ {
		if grid[row][col] == val {
			return false
		}
	}
	return true
}

func boxOK(row, col, val int, grid [][]int) bool {

	boxRow := row / 3
	boxCol := col / 3

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if grid[(boxRow*3)+x][(boxCol*3)+y] == val {
				return false
			}
		}
	}
	return true
}

// print grid
func printSol(grid [][]int) string {
	ret := ""
	for row := 0; row < 9; row++ {
		str := "| "
		for col := 0; col < 9; col++ {
			str += strconv.Itoa(grid[row][col]) + " | "
		}
		ret = ret + str + "\n"
	}
	return ret
}

// used to generate mega if block so i don't screw up by hand
func parseEq(line string) {
	// chop off end
	end := line[len(line)-5:]
	line = line[0 : len(line)-5]
	arr := strings.Split(line, " + ")
	for _, pair := range arr {
		a := int(pair[0] - 'A')
		b := int(pair[1] - '1')
		fmt.Printf("grid[%d][%d] + ", a, b)
	}
	fmt.Printf("%s\n", end)
}

/*
	"B9 + B8 + C1 + H4 + H4 = 23",
	"A5 + D7 + I5 + G8 + B3 + A5 = 19",
	"I2 + I3 + F2 + E9 = 15",
	"I7 + H8 + C2 + D9 = 26",
	"I6 + A5 + I3 + B8 + C3 = 20",
	"I7 + D9 + B6 + A8 + A3 + C4 = 27",
	"C7 + H9 + I7 + B2 + H8 + G3 = 31",
	"D3 + I8 + A4 + I6 = 27",
	"F5 + B8 + F8 + I7 + F1 = 33",
	"A2 + A8 + D7 + E4 = 21",
	"C1 + I4 + C2 + I1 + A4 = 20",
	"F8 + C1 + F6 + D3 + B6 = 25",

*/
// must be checked after value has been placed
// fail fast if we ever exceed a value that needs to be present at the end
func equationOK(grid [][]int) bool {
	if grid[1][8]+grid[1][7]+grid[2][0]+grid[7][3]+grid[7][3] > 23 ||
		grid[0][4]+grid[3][6]+grid[8][4]+grid[6][7]+grid[1][2]+grid[0][4] > 19 ||
		grid[8][1]+grid[8][2]+grid[5][1]+grid[4][8] > 15 ||
		grid[8][6]+grid[7][7]+grid[2][1]+grid[3][8] > 26 ||
		grid[8][5]+grid[0][4]+grid[8][2]+grid[1][7]+grid[2][2] > 20 ||
		grid[8][6]+grid[3][8]+grid[1][5]+grid[0][7]+grid[0][2]+grid[2][3] > 27 ||
		grid[2][6]+grid[7][8]+grid[8][6]+grid[1][1]+grid[7][7]+grid[6][2] > 31 ||
		grid[3][2]+grid[8][7]+grid[0][3]+grid[8][5] > 27 ||
		grid[5][4]+grid[1][7]+grid[5][7]+grid[8][6]+grid[5][0] > 33 ||
		grid[0][1]+grid[0][7]+grid[3][6]+grid[4][3] > 21 ||
		grid[2][0]+grid[8][3]+grid[2][1]+grid[8][0]+grid[0][3] > 20 ||
		grid[5][7]+grid[2][0]+grid[5][5]+grid[3][2]+grid[1][5] > 25 {
		return false
	}
	return true
}

// final equality check
func equationValidate(grid [][]int) bool {

	if grid[1][8]+grid[1][7]+grid[2][0]+grid[7][3]+grid[7][3] != 23 ||
		grid[0][4]+grid[3][6]+grid[8][4]+grid[6][7]+grid[1][2]+grid[0][4] != 19 ||
		grid[8][1]+grid[8][2]+grid[5][1]+grid[4][8] != 15 ||
		grid[8][6]+grid[7][7]+grid[2][1]+grid[3][8] != 26 ||
		grid[8][5]+grid[0][4]+grid[8][2]+grid[1][7]+grid[2][2] != 20 ||
		grid[8][6]+grid[3][8]+grid[1][5]+grid[0][7]+grid[0][2]+grid[2][3] != 27 ||
		grid[2][6]+grid[7][8]+grid[8][6]+grid[1][1]+grid[7][7]+grid[6][2] != 31 ||
		grid[3][2]+grid[8][7]+grid[0][3]+grid[8][5] != 27 ||
		grid[5][4]+grid[1][7]+grid[5][7]+grid[8][6]+grid[5][0] != 33 ||
		grid[0][1]+grid[0][7]+grid[3][6]+grid[4][3] != 21 ||
		grid[2][0]+grid[8][3]+grid[2][1]+grid[8][0]+grid[0][3] != 20 ||
		grid[5][7]+grid[2][0]+grid[5][5]+grid[3][2]+grid[1][5] != 25 {
		return false
	}
	return true
}
