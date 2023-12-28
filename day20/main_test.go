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
			name: "simple",
			input: `broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output
`,
			want1: 11687500,
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

func TestReadModule(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 Module
	}{
		{
			name: "s1",
			args: args{
				s: "broadcaster -> a",
			},
			want: "broadcaster",
			want1: Module{
				t:     0,
				state: 0,
				links: []string{"a"},
				inps:  map[string]byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ReadModule(tt.args.s)
			if got != tt.want {
				t.Errorf("ReadModule() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReadModule() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
