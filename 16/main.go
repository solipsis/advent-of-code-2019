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

func solveB(r io.Reader) string {

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

	cpy := make([]int, len(nums))
	copy(cpy, nums)
	for x := 1; x < 10000; x++ {
		nums = append(nums, cpy...)
	}
	//fmt.Printf("Start %v\n", nums)

	//startPattern := []int{0, 1, 0, -1}

	// 77911632 too low
	//offset := 303673
	offset := 5975677
	fmt.Println(len(nums))
	fmt.Println(offset)
	fmt.Println(len(nums) - offset)
	fmt.Printf("nums: %v", nums[offset:offset+8])
	//offset := 308177

	next := make([]int, len(nums))
	for z := 0; z < 100; z++ {
		//fmt.Printf("DAVE: %v\n", nums[offset:offset+8])

		/*
			prefixSum := 0
			for _, v := range nums[offset+8:] {
				fmt.Println(v)
				prefixSum = (prefixSum + v) % 10
				fmt.Println("x:", prefixSum)
			}
			8?
			//for w := len(nums) - 1; w >= offset+8; w-- {
			//prefixSum = (prefixSum + nums[w]) % 10
			//}
		*/

		prefixSum := 0
		for w := len(nums) - 1; w >= offset; w-- {
			next[w] = (prefixSum + nums[w]) % 10
			prefixSum = next[w]
		}

		/*
			fmt.Println("z: ", z)
			fmt.Println("prefixSum: ", prefixSum)
			for i := offset; i < len(nums); i++ {
				//for idx := range nums {
				//x := 1
				//pat := expandPattern(startPattern, y+1)
				//fmt.Printf("%v\n", pat)
				updated := 0
				for j := i; j < offset+8; j++ {
					//fmt.Printf("|%d*%d |", 1, j)
					//updated = (updated + (pat[x] * v))
					//x = (x + 1) % len(pat)
					updated = (updated + nums[j]) % 10
				}
				next[i] = (updated + prefixSum) % 10
			}
		*/
		nums = next
		fmt.Printf("%v\n", nums[offset:offset+8])
	}

	//res := ""
	//for _, v := range nums[:7] {
	//res += strconv.Itoa(v)
	//}
	fmt.Printf("test %v\n", nums[:16])
	//offset, _ := strconv.Atoi(res)
	//fmt.Println("offset:", offset)
	fmt.Println("length: ", len(nums))

	fmt.Printf("%v\n", nums[offset:offset+8])
	//return nums[offset : offset+8]
	return "blah"
	//return res

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
