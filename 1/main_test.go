package main

import (
	"bytes"
	"testing"
)

func TestExampleA(t *testing.T) {

	expected := 33583
	input := "100756"
	res := solveA(bytes.NewBufferString(input))
	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
}

func TestExampleB(t *testing.T) {
	expected := 50346
	input := "100756"
	res := solveB(bytes.NewBufferString(input))
	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
}

func TestA(t *testing.T) {
	t.Logf("A: %d", solveA(open("input.txt")))
}

func TestB(t *testing.T) {
	t.Logf("B: %d", solveB(open("input.txt")))
}
