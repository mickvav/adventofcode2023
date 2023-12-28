package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestInput_Count(t *testing.T) {
	tests := []struct {
		name  string
		input string
		steps int
		want1 int
		want2 int
	}{
		{
			name:  "s0",
			steps: 6,
			input: `...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........`,
			want1: 16,
			want2: 16,
		},
		{
			name:  "s1",
			steps: 10,
			input: `...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........`,
			want1: 33,
			want2: 50,
		},
		{
			name:  "s2",
			steps: 50,
			input: `...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........`,
			want1: 42,
			want2: 1594,
		},
		{
			name:  "s2",
			steps: 100,
			input: `...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........`,
			want1: 42,
			want2: 6536,
		},
		{
			name:  "s2",
			steps: 500,
			input: `...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........`,
			want1: 42,
			want2: 167004,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := ReadInput(bufio.NewScanner(strings.NewReader(tt.input)))
			tn := time.Now()
			got1, got2 := inp.Count(tt.steps, tt.steps)
			t.Log(time.Since(tn))
			if got1 != tt.want1 {
				t.Errorf("Input.Count().1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("Input.Count().2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestExamineMap(t *testing.T) {
	type args struct {
		v           map[int]int
		expinterval int
	}
	tests := []struct {
		name    string
		args    args
		want    []seriesDesc
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			args: args{
				v: map[int]int{
					2:  10,
					6:  14,
					10: 18,
					14: 22,
				},
				expinterval: 4,
			},
			want: []seriesDesc{{
				basei:     6,
				basev:     14,
				increment: 4,
				quadratic: 0,
			}},
			wantErr: false,
		},
		{
			name: "Q",
			args: args{
				v: map[int]int{
					2:  10,
					6:  14,
					10: 19,
					14: 25,
				},
				expinterval: 4,
			},
			want: []seriesDesc{
				{
					basei:     6,
					basev:     14,
					increment: 5,
					quadratic: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExamineMap(tt.args.v, tt.args.expinterval)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExamineMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExamineMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seriesDesc_Predict(t *testing.T) {
	type fields struct {
		basei     int
		basev     int
		increment int
		quadratic int
	}
	type args struct {
		time     int
		interval int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "s",
			fields: fields{
				basei:     7,
				basev:     3,
				increment: 2,
				quadratic: 0,
			},
			args: args{
				time:     15,
				interval: 4,
			},
			want: 7,
		},
		{
			name: "q",
			fields: fields{
				basei:     1,
				basev:     2,
				increment: 1,
				quadratic: 1,
			},
			args: args{
				time:     3,
				interval: 1,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := seriesDesc{
				basei:     tt.fields.basei,
				basev:     tt.fields.basev,
				increment: tt.fields.increment,
				quadratic: tt.fields.quadratic,
			}
			if got := s.Predict(tt.args.time, tt.args.interval); got != tt.want {
				t.Errorf("seriesDesc.Predict() = %v, want %v", got, tt.want)
			}
		})
	}
}
