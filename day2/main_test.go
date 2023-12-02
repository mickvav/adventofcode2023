package main

import (
	"strings"
	"testing"
)

func TestReadlines(t *testing.T) {
	res := Readlines(
		strings.NewReader(
			"Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green"))
	if res != 1 {
		t.Fail()
	}
}
