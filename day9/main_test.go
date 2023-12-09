package main

import (
	"bufio"
	"strings"
	"testing"
)

var testInput string = `0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45`


func TestInput_Count(t *testing.T) {
	tests := []struct {
		name   string
		want1   int
		want2  int
	}{
		{
			name:   "simple",
			want1:   114,
			want2: 2,
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
