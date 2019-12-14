package main

import (
	"bytes"
	"testing"
)

func TestSimple(t *testing.T) {

	in := `10 ORE => 10 A
1 ORE => 1 B
7 A, 1 B => 1 C
7 A, 1 C => 1 D
7 A, 1 D => 1 E
7 A, 1 E => 1 FUEL`
	expect := 31

	res := solveA(bytes.NewBufferString(in))
	if res != expect {
		t.Errorf("Got %d, expected %d", res, expect)
	}
}

func TestLarger(t *testing.T) {

	in := `9 ORE => 2 A
8 ORE => 3 B
7 ORE => 5 C
3 A, 4 B => 1 AB
5 B, 7 C => 1 BC
4 C, 1 A => 1 CA
2 AB, 3 BC, 4 CA => 1 FUEL`
	expect := 165

	res := solveA(bytes.NewBufferString(in))
	if res != expect {
		t.Errorf("Got %d, expected %d", res, expect)
	}
}
