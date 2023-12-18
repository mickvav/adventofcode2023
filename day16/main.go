package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Input struct {
	m      []string
	marks  map[pos]bool
	enters map[pos]map[dir]bool
	exits  map[pos]map[dir]bool
}

type pos dir
type Beams map[pos]map[dir]bool

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:      []string{},
		marks:  map[pos]bool{},
		enters: map[pos]map[dir]bool{},
		exits:  map[pos]map[dir]bool{},
	}
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		row++
	}
	return res
}

type dir struct {
	row, col int
}

var (
	N = dir{row: -1, col: 0}
	S = dir{row: 1, col: 0}
	E = dir{row: 0, col: 1}
	W = dir{row: 0, col: -1}
)

var ltrs = map[byte]map[dir][]dir{
	'-': {
		N: {E, W},
		S: {E, W},
		E: {E},
		W: {W},
	},
	'|': {
		N: {N},
		S: {S},
		E: {N, S},
		W: {N, S},
	},
	'\\': {
		N: {W},
		S: {E},
		E: {S},
		W: {N},
	},
	'/': {
		N: {E},
		S: {W},
		E: {N},
		W: {S},
	},
	'.': {
		N: {N},
		S: {S},
		E: {E},
		W: {W},
	},
}

func (a pos) Add(b dir) pos {
	return pos{row: a.row + b.row, col: a.col + b.col}
}

func (a dir) Opposed(b dir) bool {
	return a.row == -b.row && a.col == -b.col
}

func (a dir) Reverse() dir {
	return dir{row: -a.row, col: -a.col}
}

func (inp Input) Has(point pos) (byte, bool) {
	if point.row < 0 || point.col < 0 {
		return 0, false
	}
	if point.row >= len(inp.m) {
		return 0, false
	}
	if point.col >= len(inp.m[point.row]) {
		return 0, false
	}
	return inp.m[point.row][point.col], true
}

func (inp *Input) mark(p pos) {
	inp.marks[p] = true
}

func (inp *Input) Enter(np pos, d dir) {
	if _, ok := inp.enters[np]; !ok {
		inp.enters[np] = map[dir]bool{d: true}
		return
	}
	inp.enters[np][d] = true
}

func (inp *Input) Exit(np pos, d dir) {
	if _, ok := inp.exits[np]; !ok {
		inp.exits[np] = map[dir]bool{d: true}
		return
	}
	inp.exits[np][d] = true
}

func (inp *Input) step(b Beams) (bool, Beams) {
	res := false
	new := Beams{}
	for p, v := range b {
		for d, _ := range v {
			np := p.Add(d)
			if l, ok := inp.Has(np); ok {
				if _, ok := inp.enters[np][d]; !ok {
					inp.Enter(np, d)
					for _, d2 := range ltrs[l][d] {
						if _, ok := new[np]; !ok {
							new[np] = map[dir]bool{d2: true}
						} else {
							new[np][d2] = true
						}
						res = true
						inp.Exit(np, d2)
					}
				}
				inp.mark(np)
			}
		}
	}
	return res, new
}

func (inp Input) findFirstSteps(start pos, d dir) Beams {
	b := Beams{
		start.Add(d.Reverse()): {
			d: true,
		},
	}
	_, b1 := inp.step(b)
	inp.mark(start)
	return b1
}

func (inp *Input) Clear() {
	inp.enters = map[pos]map[dir]bool{}
	inp.exits = map[pos]map[dir]bool{}
	inp.marks = map[pos]bool{}
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0

	b := inp.findFirstSteps(pos{row: 0, col: 0}, E)
	ok := true
	for ok {
		ok, b = inp.step(b)
	}
	s1 = len(inp.marks)

	for i := range inp.m {
		b := inp.findFirstSteps(pos{row: i, col: 0}, E)
		ok := true
		for ok {
			ok, b = inp.step(b)
		}
		if len(inp.marks) > s2 {
			s2 = len(inp.marks)
		}
		inp.Clear()
	}

	for i := range inp.m {
		b := inp.findFirstSteps(pos{row: i, col: len(inp.m[0]) - 1}, W)
		ok := true
		for ok {
			ok, b = inp.step(b)
		}
		if len(inp.marks) > s2 {
			s2 = len(inp.marks)
		}
		inp.Clear()
	}

	for i := range inp.m[0] {
		b := inp.findFirstSteps(pos{row: 0, col: i}, S)
		ok := true
		for ok {
			ok, b = inp.step(b)
		}
		if len(inp.marks) > s2 {
			s2 = len(inp.marks)
		}
		inp.Clear()
	}

	for i := range inp.m[0] {
		b := inp.findFirstSteps(pos{row: len(inp.m) - 1, col: i}, N)
		ok := true
		for ok {
			ok, b = inp.step(b)
		}
		if len(inp.marks) > s2 {
			s2 = len(inp.marks)
		}
		inp.Clear()
	}

	return s1, s2
}

func Readlines(file io.Reader) (int, int) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	return inp.Count()
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Ups")
	}
	defer file.Close()
	f, f1 := Readlines(file)
	fmt.Println(strconv.Itoa(f), f1)
}
