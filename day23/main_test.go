package main

import (
	"bufio"
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
			name: "simple",
			input: `#.#####################
#.......#########...###
#######.#########.#.###
###.....#.>.>.###.#.###
###v#####.#v#.###.#.###
###.>...#.#.#.....#...#
###v###.#.#.#########.#
###...#.#.#.......#...#
#####.#.#.#######.#.###
#.....#.#.#.......#...#
#.#####.#.#.#########v#
#.#...#...#...###...>.#
#.#.#v#######v###.###v#
#...#.>.#...>.>.#.###.#
#####v#.#.###v#.#.###.#
#.....#...#...#.#.#...#
#.#########.###.#.#.###
#...###...#...#...#.###
###.###.#.###v#####v###
#...#...#.#.>.>.#.>.###
#.###.###.#.###.#.#v###
#.....###...###...#...#
#####################.#`,
			want1: 94,
			want2: 0,
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

func TestInput_ScanGraph(t *testing.T) {
	type fields struct {
		m     []string
		v     map[pos]*vertex
		e     []*edge
		start pos
		end   pos
	}
	type args struct {
		v *vertex
		d dir
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantedges int
	}{
		{
			name: "s1",
			fields: fields{
				m: []string{
					"#.##",
					"#..#",
					"##.#",
				},
				v: map[pos]*vertex{{row: 0, col: 1}: {
					p: pos{row: 0, col: 1},
					e: map[dir]*edge{},
					t: vtStart,
					d: dir{},
				},
					{
						row: 2,
						col: 2,
					}: {
						p: pos{
							row: 2,
							col: 2,
						},
						e: map[dir]*edge{},
						t: vtEnd,
						d: dir{},
					}},
				e:     []*edge{},
				start: pos{row: 0, col: 1},
				end:   pos{row: 2, col: 2},
			},
			args: args{
				v: &vertex{},
				d: dir{},
			},
			wantedges: 1,
		},
		{
			name: "s2",
			fields: fields{
				m: []string{
					"#.###",
					"#.>.#",
					"###.#",
				},
				v: map[pos]*vertex{{row: 0, col: 1}: {
					p: pos{row: 0, col: 1},
					e: map[dir]*edge{},
					t: vtStart,
					d: dir{},
				},
					{
						row: 2,
						col: 3,
					}: {
						p: pos{
							row: 2,
							col: 3,
						},
						e: map[dir]*edge{},
						t: vtEnd,
						d: dir{},
					}},
				e:     []*edge{},
				start: pos{row: 0, col: 1},
				end:   pos{row: 2, col: 3},
			},
			args: args{
				v: &vertex{},
				d: dir{},
			},
			wantedges: 2,
		},
		{
			name: "s3",
			fields: fields{
				m: []string{
					"#.###",
					"#...#",
					"#...#",
					"###.#",
				},
				v: map[pos]*vertex{{row: 0, col: 1}: {
					p: pos{row: 0, col: 1},
					e: map[dir]*edge{},
					t: vtStart,
					d: dir{},
				},
					{
						row: 3,
						col: 3,
					}: {
						p: pos{
							row: 3,
							col: 3,
						},
						e: map[dir]*edge{},
						t: vtEnd,
						d: dir{},
					}},
				e:     []*edge{},
				start: pos{row: 0, col: 1},
				end:   pos{row: 3, col: 3},
			},
			args: args{
				v: &vertex{},
				d: dir{},
			},
			wantedges: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := &Input{
				m:     tt.fields.m,
				v:     tt.fields.v,
				e:     tt.fields.e,
				start: tt.fields.start,
				end:   tt.fields.end,
			}
			inp.ScanGraph(tt.fields.v[tt.fields.start], S)
			if len(inp.e) != tt.wantedges {
				t.Errorf("ScanGraph edges outcome mismatch: want %d got %d", tt.wantedges, len(inp.e))
			}
		})
	}
}
