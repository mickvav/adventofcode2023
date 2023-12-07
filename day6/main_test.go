package main

import (
	"strings"
	"testing"
)

var testInput string = `Time:      7  15   30
Distance:  9  40  200`

func TestReadlines(t *testing.T) {
	res1, res2 := Readlines(
		strings.NewReader(testInput),
	)
	if res1 != 288 {
		t.Fail()
	}
	if res2 != 71503 {
		t.Fail()
	}
}
