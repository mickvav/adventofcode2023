package main

import (
	"bufio"
	"strings"
	"testing"
)

var testInput string = `LLR

AAA = (BBB, BBB)
BBB = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)`

var testInput2 string = `LR

11A = (11B, XXX)
11B = (XXX, 11Z)
11Z = (11B, XXX)
22A = (22B, XXX)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22B, 22B)
AAA = (BBB, BBB)
BBB = (CCC, CCC)
CCC = (DDD, DDD)
DDD = (ZZZ, ZZZ)
ZZZ = (ZZZ, ZZZ)
XXX = (XXX, XXX)`


func TestInput_Count2(t *testing.T) {
	type fields struct {
		instructions string
		L            map[string]string
		R            map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name:   "simple",
			fields: fields{},
			want:   12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := ReadInput(bufio.NewScanner(strings.NewReader(testInput2)))
			if got := inp.Count2(); got != tt.want {
				t.Errorf("Input.Count2() = %v, want %v", got, tt.want)
			}
		})
	}
}
