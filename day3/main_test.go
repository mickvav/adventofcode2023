package main

import (
	"strings"
	"testing"
)

func TestReadlines(t *testing.T) {
	res := Readlines(
		strings.NewReader(
			"467..114..\n" +
				"...*......\n" +
				"..35..633.\n" +
				"......#...\n" +
				"617*......\n" +
				".....+.58.\n" +
				"..592.....\n" +
				"......755.\n" +
				"...$.*....\n" +
				".664.598..",
		))
	if res.Process() != 4361 {
		t.Fail()
	}
	if res.FindAllGears() != 467835 {
		t.Fail()
	}
}
