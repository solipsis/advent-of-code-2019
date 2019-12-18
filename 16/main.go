package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func solveA(r io.Reader) string {
	sc := bufio.NewScanner(r)

	nums := make([]int, 0)
	for sc.Scan() {
		arr := strings.Split(sc.Text(), "")
		for _, s := range arr {
			i, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			nums = append(nums, i)
		}
	}

	startPattern := []int{0, 1, 0, -1}

	next := make([]int, len(nums))
	for z := 0; z < 100; z++ {
		for idx := range nums {
			x := 1
			pat := expandPattern(startPattern, idx+1)
			//fmt.Printf("%v\n", pat)
			updated := 0
			for _, v := range nums {
				//fmt.Printf("|%d*%d |", pat[x], v)
				updated += pat[x] * v
				x = (x + 1) % len(pat)
			}
			next[idx] = abs(updated) % 10
		}
		nums = next
		//fmt.Printf("%v\n", nums)
	}

	res := ""
	for _, v := range nums[:8] {
		res += strconv.Itoa(v)
	}

	return res
}

func solveB(r io.Reader) int {

	sc := bufio.NewScanner(r)

	nums := make([]int, 0)
	for sc.Scan() {
		arr := strings.Split(sc.Text(), "")
		for _, s := range arr {
			i, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			nums = append(nums, i)
		}
	}

	// duplicate input 10k times
	cpy := make([]int, len(nums))
	copy(cpy, nums)
	for x := 1; x < 10000; x++ {
		nums = append(nums, cpy...)
	}

	//offset := 303673
	offset := 5975677

	// the first n-1 values will all be zeroes
	// and the coefficient for all values after will be 1

	next := make([]int, len(nums))
	for z := 0; z < 100; z++ {
		prefixSum := 0
		// iterate in reverse to avoid duplicate computation
		for w := len(nums) - 1; w >= offset; w-- {
			next[w] = (prefixSum + nums[w]) % 10
			prefixSum = next[w]
		}
		nums = next
	}

	// turn list into number
	res := 0
	for _, v := range nums[offset : offset+8] {
		res += v
		res *= 10
	}

	return res / 10

}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func expandPattern(pat []int, n int) []int {

	res := make([]int, 0)
	for _, v := range pat {
		for x := 0; x < n; x++ {
			res = append(res, v)
		}
	}

	return res

	//return append(res[1:], res[0])
}

func main() {
	input := open("input.txt")
	//fmt.Printf("A: %s\n", solveA(input))
	//input = open("input.txt")
	fmt.Printf("B: %d\n", solveB(input))
}

func open(fname string) io.Reader {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	return f
}
