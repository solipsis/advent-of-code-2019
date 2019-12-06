package main

import (
	"fmt"
	"strconv"
	"strings"
)

var grid [9][9]int
var preset = make(map[int]bool)
var debug int

func main() {

	// apply orignal constraints
	grid[0][8] = 1
	grid[1][1] = 1
	grid[1][2] = 2
	grid[2][6] = 2
	grid[3][8] = 2
	grid[4][1] = 2
	grid[6][6] = 1
	grid[6][7] = 2
	grid[7][0] = 1
	grid[7][5] = 2
	grid[8][3] = 1
	grid[0][0] = 9
	grid[0][1] = 4
	grid[0][2] = 6
	grid[0][3] = 7
	grid[0][4] = 5
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
	preset[83] = true

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

	//printSol()
	solve(0, 4)

}

func solve(r, c int) {
	debug++
	if debug%1000000 == 0 {
		fmt.Println("\n\n----------------------------------------")
		printSol()
	}
	//fmt.Println(r, c)

	// recursive base case
	if r == 9 {
		if equationValidate() {
			fmt.Println("correct")
			printSol()
			panic("done")
			return
		} else {
			fmt.Println("final validation failed")
			return
		}
	}

	// is this square preset? if so skip
	if preset[r*10+c] {
		if c == 8 {
			solve(r+1, 0)
		} else {
			solve(r, c+1)
		}
		return
	}

	// place values 1-9 in this square
	for v := 1; v <= 9; v++ {
		if rowOK(r, v) && colOK(c, v) && boxOK(r, c, v) {
			grid[r][c] = v
			if !equationOK() {
				//grid[r][c] = 0
				continue
			}

			// this is valid, move to the next column, or row if this is the last column
			if c == 8 {
				solve(r+1, 0)
			} else {
				solve(r, c+1)
			}
		}
	}
	grid[r][c] = 0

}

func rowOK(row, val int) bool {
	for col := 0; col < 9; col++ {
		if grid[row][col] == val {
			//fmt.Println("row fail")
			return false
		}
	}
	return true
}

func colOK(col, val int) bool {
	for row := 0; row < 9; row++ {
		if grid[row][col] == val {
			//fmt.Println("col fail")
			return false
		}
	}
	return true
}

func boxOK(row, col, val int) bool {

	boxRow := row / 3
	boxCol := col / 3

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if grid[(boxRow*3)+x][(boxCol*3)+y] == val {
				//fmt.Println("box fail")
				return false
			}
		}
	}
	return true
}

func printSol() {
	for row := 0; row < 9; row++ {
		str := "| "
		for col := 0; col < 9; col++ {
			str += strconv.Itoa(grid[row][col]) + " | "
		}
		fmt.Println(str)
	}
}

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

// must be checked after value has been placed
func equationOK() bool {
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

func equationValidate() bool {

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
