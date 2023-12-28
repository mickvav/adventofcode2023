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
		want2 int
	}{
		{
			name: "s0",
			input: `111111111111
999999999991
999999999991
999999999991
999999999991`,
			want1: 59,
			want2: 71,
		},
		{
			name: "simple",
			input: `2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533`,
			want1: 102,
			want2: 94,
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

func Test_state_Neighbours2(t *testing.T) {
	type fields struct {
		row          int
		col          int
		prevdir      byte
		prevdistance byte
	}
	type args struct {
		inp Input
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []state
	}{
		{
			name: "East",
			fields: fields{
				row:          0,
				col:          0,
				prevdir:      'E',
				prevdistance: 0,
			},
			args: args{
				inp: Input{
					m:      []string{"11111", "99999"},
					v:      [][]int{{1, 1, 1, 1, 1}, {9, 9, 9, 9, 9}},
					marks:  map[state]bool{},
					values: map[state]int{},
				},
			},
			want: []state{{
				row:          0,
				col:          1,
				prevdir:      'E',
				prevdistance: 1,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state{
				row:          tt.fields.row,
				col:          tt.fields.col,
				prevdir:      tt.fields.prevdir,
				prevdistance: tt.fields.prevdistance,
			}
			if got := s.Neighbours2(tt.args.inp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("state.Neighbours2() = %v, want %v", got, tt.want)
			}
		})
	}
}
