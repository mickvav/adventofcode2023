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
		want2 uint64
	}{
		{
			name:  "sx2",
			input: "?###???????? 3,2,1",
			want1: 10,
			want2: 506250,
		},
		{
			name:  "s2",
			input: "?? 1",
			want1: 2,
			want2: 252,
		},
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
			name:  "s1",
			input: "????.######..#####. 1,6,5",
			want1: 4,
			want2: 2500,
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
		{
			name:  "....",
			input: `??????? 1,1,1`,
			want1: 10,
			want2: 3268760,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := ReadInput(bufio.NewScanner(strings.NewReader(tt.input)))
			got1, got2 := inp.Count()
			if got1 != tt.want1 {
				t.Errorf("Input.Count().1 = %v, want %v", got1, tt.want1)
				inp.mh = mMatcher{
					h: map[mask]map[intStructure]uint16{},
				}
				inp.Count()
			}
			if got2 != tt.want2 {
				t.Errorf("Input.Count().2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_intStructure_Concat(t *testing.T) {
	type args struct {
		b intStructure
	}
	tests := []struct {
		name string
		is   intStructure
		args args
		want intStructure
	}{
		{
			name: "s0",
			is:   intStructure([]byte{1, 2, 3}),
			args: args{
				b: intStructure([]byte{0, 2, 3}),
			},
			want: intStructure([]byte{1, 2, 3, 2, 3}),
		},
		{
			name: "s1",
			is:   intStructure([]byte{1, 2, 3}),
			args: args{
				b: intStructure([]byte{2, 3}),
			},
			want: intStructure([]byte{1, 2, 5, 3}),
		},
		{
			name: "s2",
			is:   intStructure([]byte{0}),
			args: args{
				b: intStructure([]byte{0, 1, 3}),
			},
			want: intStructure([]byte{0, 1, 3}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.is.Concat(tt.args.b); got != tt.want {
				t.Errorf("intStructure.Concat() = %v, want %v", got.String(), tt.want.String())
				tt.is.Concat(tt.args.b)
			}
		})
	}
}

func Test_guess_IntervalStructure(t *testing.T) {
	tests := []struct {
		name string
		g    guess
		want intStructure
	}{
		// TODO: Add test cases.
		{
			name: "s0",
			g:    0,
			want: intStructure([]byte{0, 0}),
		},
		{
			name: "s1",
			g:    0b010,
			want: intStructure([]byte{0, 1, 0}),
		},
		{
			name: "s2",
			g:    0b0110,
			want: intStructure([]byte{0, 2, 0}),
		},
		{
			name: "s3",
			g:    0b0101,
			want: intStructure([]byte{1, 1, 0}),
		},
		{
			name: "s4",
			g:    0b1000000000000111,
			// .... FEDCBA9876543210
			want: intStructure([]byte{3, 1}),
		},
		{
			name: "s5",
			g:    0b0100000000000111,
			// .... FEDCBA9876543210
			want: intStructure([]byte{3, 1, 0}),
		},
		{
			name: "s6",
			g:    0b0100000000000110,
			// .... FEDCBA9876543210
			want: intStructure([]byte{0, 2, 1, 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.calcIntervalStructure(); got != tt.want {
				t.Errorf("guess.IntervalStructure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intStructure_IsStartOf(t *testing.T) {
	type args struct {
		m []int
	}
	tests := []struct {
		name string
		is   intStructure
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "s0",
			is:   intStructure([]byte{0, 1, 2}),
			args: args{
				m: []int{1, 2, 3, 4},
			},
			want: true,
		},
		{
			name: "s1",
			is:   intStructure([]byte{0, 1, 2}),
			args: args{
				m: []int{1, 3, 3, 4},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.is.IsStartOf(tt.args.m); got != tt.want {
				t.Errorf("intStructure.IsStartOf() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
		{
			name: "simple",
			args: args{
				s: ".#??#..........###...?.... 1,2,3",
				//  0123456789ABCDEF0123456789
			},
			want: Line{
				orig:   ".#??#..........###...?....",
				length: 10 + 16,
				masks: []mask{
					{
						unknown: 0b0000000000001100,
						//             FEDCBA9876543210
						damaged:     0b1000000000010010,
						operational: 0b0111111111100001,
					}, {
						unknown:     0b0000000000100000,
						damaged:     0b0000000000000011,
						operational: 0b1111111111011100,
					},
				},
				checksums: []int{1, 2, 3},
			},
		},
		{
			name: "32",
			args: args{
				s: ".#.............#..#............# 1,1,1,1",
				//  0123456789ABCDEF0123456789ABCDEF
			},
			want: Line{
				orig:   ".#.............#..#............#",
				length: 32,
				masks: []mask{{
					unknown: 0b0,
					//             FEDCBA9876543210
					damaged:     0b1000000000000010,
					operational: 0b0111111111111101,
				},
					{
						unknown:     0b0,
						damaged:     0b1000000000000100,
						operational: 0b0111111111111011,
					},
				},
				checksums: []int{1, 1, 1, 1},
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

func Test_mMatcher_Mask(t *testing.T) {
	type fields struct {
		h map[mask]map[intStructure]uint16
	}
	type args struct {
		msk mask
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[intStructure]uint16
	}{
		// TODO: Add test cases.
		{
			name: "simple",
			fields: fields{
				h: map[mask]map[intStructure]uint16{},
			},
			args: args{
				msk: mask{
					unknown:     0b0000000000000001,
					damaged:     0b0000000000000000,
					operational: 0b1111111111111110,
				},
			},
			want: map[intStructure]uint16{
				intStructure([]byte{0, 0}): 1,
				intStructure([]byte{1, 0}): 1,
			},
		},
		{
			name: "simple",
			fields: fields{
				h: map[mask]map[intStructure]uint16{},
			},
			args: args{
				msk: mask{
					unknown:     0b0000000000000111,
					damaged:     0b0000000000000000,
					operational: 0b1111111111111000,
				},
			},
			want: map[intStructure]uint16{
				intStructure([]byte{0, 0}):    1,
				intStructure([]byte{1, 0}):    1,
				intStructure([]byte{0, 1, 0}): 2,
				intStructure([]byte{1, 1, 0}): 1,
				intStructure([]byte{2, 0}):    1,
				intStructure([]byte{3, 0}):    1,
				intStructure([]byte{0, 2, 0}): 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mMatcher{
				h: tt.fields.h,
			}
			if got := m.Mask(tt.args.msk); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mMatcher.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
