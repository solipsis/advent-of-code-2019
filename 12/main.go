package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type body struct {
	x, y, z    int
	vx, vy, vz int
}

func solveA(r io.Reader) int {
	sc := bufio.NewScanner(r)

	var bodies []body
	for sc.Scan() {
		// ulitimate lazy input parsing
		rep := strings.NewReplacer(" ", "", "<", "", ">", "", "=", "", "x", "", "y", "", "z", "")
		text := rep.Replace(sc.Text())
		arr := strings.Split(text, ",")

		x, _ := strconv.Atoi(arr[0])
		y, _ := strconv.Atoi(arr[1])
		z, _ := strconv.Atoi(arr[2])

		bodies = append(bodies, body{x: x, y: y, z: z})
	}

	for step := 1; step <= 1000; step++ {
		applyGravity(bodies)
		applyVelocity(bodies)
		//for _, b := range bodies {
		//fmt.Printf("x: %d, y: %d, z: %d, vx: %d, vy: %d, vz: %d\n", b.x, b.y, b.z, b.vx, b.vy, b.vz)
		//}
		//fmt.Println(step, "*********************************************")
	}

	// total energy
	energy := 0
	for _, b := range bodies {
		pe := abs(b.x) + abs(b.y) + abs(b.z)
		ke := abs(b.vx) + abs(b.vy) + abs(b.vz)
		energy += pe * ke
	}

	return energy
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func applyVelocity(bodies []body) {
	for i := 0; i < len(bodies); i++ {
		bodies[i].x += bodies[i].vx
		bodies[i].y += bodies[i].vy
		bodies[i].z += bodies[i].vz
	}
}

func applyGravity(bodies []body) {
	cpy := make([]body, len(bodies))
	for i, b := range bodies {
		cpy[i] = body{
			x:  b.x,
			y:  b.y,
			z:  b.z,
			vx: b.vx,
			vy: b.vy,
			vz: b.vz,
		}
	}

	// for every pair
	for i := 0; i < len(bodies); i++ {
		for j := i + 1; j < len(bodies); j++ {
			// X
			if cpy[i].x < cpy[j].x {
				bodies[i].vx++
				bodies[j].vx--
			} else if cpy[i].x > cpy[j].x {
				bodies[i].vx--
				bodies[j].vx++
			}
			// Y
			if cpy[i].y < cpy[j].y {
				bodies[i].vy++
				bodies[j].vy--
			} else if cpy[i].y > cpy[j].y {
				bodies[i].vy--
				bodies[j].vy++
			}
			// z
			if cpy[i].z < cpy[j].z {
				bodies[i].vz++
				bodies[j].vz--
			} else if cpy[i].z > cpy[j].z {
				bodies[i].vz--
				bodies[j].vz++
			}
		}
	}

}

func solveB(r io.Reader) int {
	return -1
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
