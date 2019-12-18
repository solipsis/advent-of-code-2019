package main

import (
	"bytes"
	"testing"
)

func TestExampleA(t *testing.T) {

	var tests = []struct {
		input, expect string
	}{
		{"80871224585914546619083218645595", "24176176"},
		{"19617804207202209144916044189917", "73745418"},
		{"69317163492948606335995924319873", "52432133"},
	}

	for _, tt := range tests {
		res := solveA(bytes.NewBufferString(tt.input))
		if res != tt.expect {
			t.Errorf("Got %s, expected %s", res, tt.expect)
		}
	}
}
