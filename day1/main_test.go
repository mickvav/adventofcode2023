package main

import (
	"strings"
	"testing"
)

func TestReadlines(t *testing.T) {
	res := Readlines(
		strings.NewReader(
			"pdrss6oneone4fournine"))
	if res != 64 {
		t.Fail()
	}
}
