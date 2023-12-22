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
		want1 int
		want2 int
	}{
		{
			name: "s0",
			input: `1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9
`,
			want1: 5,
			want2: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := ReadInput(bufio.NewScanner(strings.NewReader(tt.input)))
			tn := time.Now()
			got1, got2 := inp.Count()
			inp.Print()
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

func TestReadBrick(t *testing.T) {
	type args struct {
		s  string
		rn int
	}
	tests := []struct {
		name string
		args args
		want brick
	}{
		// TODO: Add test cases.
		{
			name: "low",
			args: args{
				s:  "1,0,1~1,2,1",
				rn: 1,
			},
			want: brick{
				rn: 1,
				img: [4]uint64{
					0b10 | (0b10 << 16) | (0b10 << 32),
					0,
				},
				z: 1,
				h: 1,
			},
		},
		{
			name: "high",
			args: args{
				s:  "0,8,2~0,9,2",
				rn: 1,
			},
			want: brick{
				rn: 1,
				img: [4]uint64{0, 0,
					0b1 | 0b1<<16,
				},
				z: 2,
				h: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadBrick(tt.args.s, tt.args.rn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadBrick() = %v, want %v", got, tt.want)
			}
		})
	}
}
