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
			name: "simple",
			input: `19, 13, 30 @ -2,  1, -2
18, 19, 22 @ -1, -1, -2
20, 25, 34 @ -2, -2, -4
12, 31, 28 @ -1, -2, -1
20, 19, 15 @  1, -5, -3`,
			want1: 2,
			want2: 47,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minCoord = 7
			pfactor = 1.0
			maxCoord = 27
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

func TestHail_Intersect2d(t *testing.T) {
	type fields struct {
		p pos
		d dir
	}
	type args struct {
		h2 Hail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    pos
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "s0",
			fields: fields{
				p: pos{
					x: 20,
					y: 19,
					z: 0,
				},
				d: dir{
					x: 0,
					y: -5,
					z: 0,
				},
			},
			args: args{
				h2: Hail{
					p: pos{
						x: 20,
						y: 25,
						z: 0,
					},
					d: dir{
						x: -2,
						y: -2,
						z: 0,
					},
				},
			},
			want:    pos{},
			wantErr: true,
		},
		{
			name: "s1",
			fields: fields{
				p: pos{
					x: 20,
					y: 25,
					z: 0,
				},
				d: dir{
					x: -2,
					y: -2,
					z: 0,
				},
			},
			args: args{
				h2: Hail{
					p: pos{
						x: 20,
						y: 19,
						z: 0,
					},
					d: dir{
						x: 0,
						y: -5,
						z: 0,
					},
				},
			},
			want:    pos{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h1 := Hail{
				p: tt.fields.p,
				d: tt.fields.d,
			}
			minCoord = 7
			maxCoord = 27
			got, err := h1.Intersect2d(tt.args.h2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hail.Intersect2d() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hail.Intersect2d() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "simple",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func TestHail_GoesThroughBbox2D(t *testing.T) {
	type fields struct {
		p pos
		d dir
	}
	tests := []struct {
		name   string
		fields fields
		minc   float64
		maxc   float64
		want   bool
	}{
		// TODO: Add test cases.
		{
			name:   "s1",
			fields: fields{p: pos{x: 0, y: 0, z: 0}, d: dir{x: 1, y: 1, z: 0}},
			minc:   10,
			maxc:   20,
			want:   true,
		},
		{
			name: "s2",
			fields: fields{
				p: pos{
					x: 0,
					y: 10,
					z: 0,
				},
				d: dir{
					x: 10,
					y: 1,
					z: 0,
				},
			},
			minc: 10,
			maxc: 20,
			want: true,
		},
		{
			name: "s3",
			fields: fields{
				p: pos{
					x: 10,
					y: 0,
					z: 0,
				},
				d: dir{
					x: 1,
					y: 10,
					z: 0,
				},
			},
			minc: 10,
			maxc: 20,
			want: true,
		},
		{
			name: "s4",
			fields: fields{
				p: pos{
					x: 0,
					y: 0,
					z: 0,
				},
				d: dir{
					x: 1,
					y: 10,
					z: 0,
				},
			},
			minc: 10,
			maxc: 20,
			want: false,
		},
		{
			name: "s5",
			fields: fields{
				p: pos{
					x: 30,
					y: 10,
					z: 0,
				},
				d: dir{
					x: -10,
					y: 1,
					z: 0,
				},
			},
			minc: 10,
			maxc: 20,
			want: true,
		},
		{
			name: "l0",
			fields: fields{
				// 233210433951170, 272655040388795, 179982504986147 @ 39, -98, 166
				p: pos{
					x: 233210433951170,
					y: 272655040388795,
					z: 179982504986147,
				},
				d: dir{x: 39, y: -98, z: 166},
			},
			minc: 200000000000000,
			maxc: 400000000000000,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h1 := Hail{
				p: tt.fields.p,
				d: tt.fields.d,
			}
			minCoord = tt.minc
			maxCoord = tt.maxc
			if got := h1.GoesThroughBbox2D(); got != tt.want {
				t.Errorf("Hail.GoesThroughBbox2D() = %v, want %v", got, tt.want)
				h1.GoesThroughBbox2D()
			}
		})
	}
}

func TestHail_MinDistance2(t *testing.T) {
	type fields struct {
		p pos
		d dir
	}
	type args struct {
		h2 Hail
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		// TODO: Add test cases.
		{
			name: "triv",
			fields: fields{
				p: pos{x: 0, y: 0, z: 0},
				d: dir{x: 1, y: 1, z: 1},
			},
			args: args{
				h2: Hail{
					p: pos{x: 1, y: 0, z: 0},
					d: dir{x: -1, y: 1, z: 1},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h1 := Hail{
				p: tt.fields.p,
				d: tt.fields.d,
			}
			if got := h1.MinDistance2(tt.args.h2); got != tt.want {
				t.Errorf("Hail.MinDistance2() = %v, want %v", got, tt.want)
			}
		})
	}
}
