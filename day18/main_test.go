package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestInput_Count(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want1 int
		want2 int64
	}{
		{
			name: "synthetic2",
			//  ##
			//  ###
			//  ###
			input: `R 2 (#000020)
U 1 (#000013)
L 1 (#000012)
U 1 (#000013)
L 1 (#000012)
D 2 (#000021)`,
			want1: 8,
			want2: 8,
		},

		{
			name: "synthetic",
			input: `R 2 (#000020)
U 2 (#000023)
L 2 (#000022)
D 2 (#000021)`,
			want1: 9,
			want2: 9,
		},
		{
			name: "simple",
			input: `R 6 (#70c710)
D 5 (#0dc571)
L 2 (#5713f0)
D 2 (#d2c081)
R 2 (#59c680)
D 2 (#411b91)
L 5 (#8ceee2)
U 2 (#caa173)
L 1 (#1b58a2)
U 2 (#caa171)
R 2 (#7807d2)
U 3 (#a77fa3)
L 2 (#015232)
U 2 (#7a21e3)`,
			want1: 62,
			want2: 952408144115,
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

func TestReadLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Line
	}{
		{
			name: "s1",
			args: args{
				"U 10 (#015232)",
			},
			want: Line{
				direction:  'U',
				length:     10,
				direction2: 'L',
				length2:    5411,
			},
		},
		{
			name: "s2",
			args: args{
				s: "R 6 (#70c710)",
			},
			want: Line{
				direction:  'R',
				length:     6,
				direction2: 'R',
				length2:    0x70c71,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadLine(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
