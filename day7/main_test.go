package main

import (
	"strings"
	"testing"
)

var testInput string = `32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`

func TestReadlines(t *testing.T) {
	res1, res2 := Readlines(
		strings.NewReader(testInput),
	)
	if res1 != 6440 {
		t.Fail()
	}
	if res2 != 5905 {
		t.Fail()
	}
}
