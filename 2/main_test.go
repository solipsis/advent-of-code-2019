package main

import (
	"bytes"
	"testing"
)

func TestExampleA(t *testing.T) {

	tt := []struct {
		input string
		out   int
	}{
		{"1,9,10,3,2,3,11,0,99,30,40,50", 3500},
		{"1,0,0,0,99", 2},
		{"2,3,0,3,99", 2},
		{"1,1,1,4,99,5,6,0,99", 30},
	}

	for _, test := range tt {
		res := solveA(bytes.NewBufferString(test.input))
		if res != test.out {
			t.Errorf("Got %d, expected %d", res, test.out)
		}
	}
}

func TestA(t *testing.T) {
	t.Logf("A: %d", solveA(open("input.txt")))
}

func TestB(t *testing.T) {
	t.Logf("B: %d", solveB(open("input.txt")))
}
