package main

import (
	"bufio"
	"strings"
	"testing"
)

var testInput string = 
`..F7.
.FJ|.
SJ.L7
|F--J
LJ...`


func TestInput_Count(t *testing.T) {
	tests := []struct {
		name   string
		want1   int
		want2  int
	}{
		{
			name:   "simple",
			want1:   8,
			want2: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := ReadInput(bufio.NewScanner(strings.NewReader(testInput)))
			got1, got2 := inp.Count(); 
			if got1 != tt.want1 {
				t.Errorf("Input.Count().1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("Input.Count().2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
