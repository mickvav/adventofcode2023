package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

var testInput string = `..F7.
.FJ|.
SJ.L7
|F--J
LJ...`

func TestInput_Count(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want1 int
		want2 int
	}{
		{
			name: "simple",
			input: `#.#.### 1,1,3
.#...#....###. 1,1,3
.#.###.#.###### 1,3,1,6
####.#...#... 4,1,1
#....######..#####. 1,6,5
.###.##....# 3,2,1`,
			want1: 6,
			want2: 6,
		},
		{
			name: "unknowns",
			input: `???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1
`,
			want1: 21,
			want2: 525152,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := ReadInput(bufio.NewScanner(strings.NewReader(tt.input)))
			got1, got2 := inp.Count()
			if got1 != tt.want1 {
				t.Errorf("Input.Count().1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("Input.Count().2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestLine_StartIteration(t *testing.T) {
	type fields struct {
		orig        string
		length      uint64
		unknowns    uint128
		damaged     uint128
		operational uint128
		checksums   []uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   LineIterator
	}{
		{
			name: "simple",
			fields: fields{
				orig:        "???.###",
				length:      7,
				unknowns:    uint128{hi: 0, lo: 0b0000111},
				damaged:     uint128{hi: 0, lo: 0b1110000},
				operational: uint128{hi: 0, lo: 0b0001000},
				checksums:   []uint64{1, 1, 3},
			},
			want: LineIterator{
				length:        7,
				damagedlength: 5,
				checksums:     []uint64{1, 1, 3},
				state:         []uint64{0, 1, 1, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				orig:        tt.fields.orig,
				length:      tt.fields.length,
				unknowns:    tt.fields.unknowns,
				damaged:     tt.fields.damaged,
				operational: tt.fields.operational,
				checksums:   tt.fields.checksums,
			}
			if got := l.StartIteration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line.StartIteration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLineIterator_Step(t *testing.T) {
	type fields struct {
		length        uint64
		damagedlength uint64
		checksums     []uint64
		state         []uint64
	}
	tests := []struct {
		name     string
		fields   fields
		want     bool
		hint     uint128
		newstate []uint64
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			fields: fields{
				length:        10,
				damagedlength: 5,
				checksums:     []uint64{2, 3},
				state:         []uint64{3, 1, 1},
			},
			want:     true,
			hint:     uint128{hi: 0, lo: 0},
			newstate: []uint64{2, 2, 1},
		},
		{
			name: "unhinted",
			fields: fields{
				length:        10,
				damagedlength: 5,
				checksums:     []uint64{2, 3},
				state:         []uint64{3, 1, 1},
				// ...##.###. -> 0b0111011000

			},
			want:     true,
			hint:     uint128{hi: 0, lo: 0},
			newstate: []uint64{2, 2, 1},
			// ..##..###. -> 0b0111001100
		},
		{
			name: "hinted",
			fields: fields{
				length:        10,
				damagedlength: 5,
				checksums:     []uint64{2, 3},
				state:         []uint64{3, 1, 1},
				// ...##.###. -> 0b0111011000

			},
			want: true,
			// ...##.##!.
			// ........!.
			hint:     uint128{hi: 0, lo: 0b0100000000},
			newstate: []uint64{2, 1, 2},
			// ..##.###.. -> 0b0011101100
		},

		{
			name: "end",
			fields: fields{
				length:        10,
				damagedlength: 5,
				checksums:     []uint64{2, 3},
				state:         []uint64{0, 1, 4},
			},
			want:     false,
			hint:     uint128{hi: 0, lo: 0},
			newstate: []uint64{0, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &LineIterator{
				length:        tt.fields.length,
				damagedlength: tt.fields.damagedlength,
				checksums:     tt.fields.checksums,
				state:         tt.fields.state,
			}
			if got := it.Step(tt.hint); got != tt.want {
				t.Errorf("LineIterator.Step() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(it.state, tt.newstate) {
				t.Errorf("LineIterator.Step() new state mismatch = %v, want %v", it.state, tt.newstate)
			}
		})
	}
}

func TestLineIterator_Repr(t *testing.T) {
	type fields struct {
		length        uint64
		damagedlength uint64
		checksums     []uint64
		state         []uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint128
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			fields: fields{
				length:        5,
				damagedlength: 3,
				checksums:     []uint64{1, 2},
				state:         []uint64{0, 1, 1},
			},
			want: uint128{hi: 0, lo: 0b01101},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &LineIterator{
				length:        tt.fields.length,
				damagedlength: tt.fields.damagedlength,
				checksums:     tt.fields.checksums,
				state:         tt.fields.state,
			}
			if got := it.Repr(); got != tt.want {
				t.Errorf("LineIterator.Repr() = %v, want %v", got, tt.want)
			}
		})
	}
}
