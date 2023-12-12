package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Input struct {
	m     []string
	marks map[dir]bool
	srow  int
	scol  int
	sltr  string
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:     []string{},
		marks: map[dir]bool{},
		srow:  0,
		scol:  0,
	}
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		for i, v := range strings.Split(line, "") {
			if v == "S" {
				res.scol = i
				res.srow = row
			}
		}
		row++
	}
	return res
}

type dir struct {
	row, col int
}

var ltrs = map[byte][]dir{
	'-': {{
		row: 0,
		col: 1,
	},
		{
			row: 0,
			col: -1,
		}},
	'|': {{
		row: 1,
		col: 0,
	},
		{
			row: -1,
			col: 0,
		}},
	'L': {{
		row: -1,
		col: 0,
	},
		{
			row: 0,
			col: 1,
		}},
	'F': {{
		row: 1,
		col: 0,
	},
		{
			row: 0,
			col: 1,
		}},
	'7': {{
		row: 1,
		col: 0,
	},
		{
			row: 0,
			col: -1,
		}},
	'J': {{
		row: -1,
		col: 0,
	},
		{
			row: 0,
			col: -1,
		}},
}

var nextdir = map[byte]map[dir]dir{}

var alldirs = []dir{
	{
		row: -1,
		col: 0,
	},
	{
		row: 1,
		col: 0,
	},
	{
		row: 0,
		col: -1,
	},
	{
		row: 0,
		col: 1,
	},
}

func (a dir) Add(b dir) dir {
	return dir{row: a.row + b.row, col: a.col + b.col}
}

func (a dir) Opposed(b dir) bool {
	return a.row == -b.row && a.col == -b.col
}

func (a dir) Reverse() dir {
	return dir{row: -a.row, col: -a.col}
}

func (inp Input) Has(point dir) (byte, bool) {
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

func (inp *Input) mark(pos dir) {
	inp.marks[pos] = true
}
func (inp Input) step(pos, angle dir) (nextpos, nextangle dir, err error) {
	if l, ok := inp.Has(pos); ok {
		nextangle = nextdir[l][angle]
		nextpos = pos.Add(nextangle)
	} else {
		return pos, angle, fmt.Errorf("no next pos %v", nextpos)
	}
	return
}

func (inp Input) findFirstSteps() (pos []dir, angles []dir) {
	start := dir{row: inp.srow, col: inp.scol}
	inp.mark(start)
	for _, d := range alldirs {
		s := d.Add(start)
		if r, ok := inp.Has(s); ok {
			if dirs, ok := ltrs[r]; ok {
				if dirs[0].Opposed(d) {
					pos = append(pos, s)
					inp.mark(s)
					angles = append(angles, d)
				}
				if dirs[1].Opposed(d) {
					pos = append(pos, s)
					inp.mark(s)
					angles = append(angles, d)
				}

			}
		}
	}
	return pos, angles
}

func fillNextDir() {
	for k, v := range ltrs {
		for i, d := range v {
			if _, ok := nextdir[k]; !ok {
				nextdir[k] = map[dir]dir{}
			}
			nextdir[k][d.Reverse()] = v[1-i]
		}
	}
}

func initialLetter(angles []dir) string {
	for k, v := range ltrs {
		if (v[0] == angles[0] && v[1] == angles[1]) || (v[0] == angles[1] && v[1] == angles[0]) {
			return string([]byte{k})
		}
	}
	return "."
}
func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	fillNextDir()
	pos, angles := inp.findFirstSteps()
	if len(pos) != 2 {
		return 0, 0
	}
	s1 = 1
	inp.sltr = initialLetter(angles)
	var e1, e2 error
	e1 = nil
	e2 = nil
	for pos[0] != pos[1] && e1 == nil && e2 == nil {
		pos[0], angles[0], e1 = inp.step(pos[0], angles[0])
		pos[1], angles[1], e2 = inp.step(pos[1], angles[1])
		inp.mark(pos[0])
		inp.mark(pos[1])
		s1 += 1
	}
	state := 0 // 0 - out, 1 - below, 2 - above, 3 - in
	for i, v := range inp.m {
		state = 0
		for i2, v2 := range strings.Split(v, "") {
			p := dir{col: i2, row: i}

			if v, ok := inp.marks[p]; (ok && v) || v2 == "S" {
				v3 := v2
				if v2 == "S" {
					v3 = inp.sltr
				}
				switch v3 {
				case "|":
					switch state {
					case 0:
						state = 3
					case 3:
						state = 0
					}

				case "F":
					switch state {
					case 0:
						state = 1
					case 3:
						state = 2
					}

				case "J":
					switch state {
					case 1:
						state = 3
					case 2:
						state = 0
					}

				case "7":
					switch state {
					case 1:
						state = 0
					case 2:
						state = 3
					}

				case "L":
					switch state {
					case 0:
						state = 2
					case 3:
						state = 1
					}
				}
			}
			if _, ok := inp.marks[p]; !ok && state == 3 {
				s2 += 1
				//println(p.row, p.col)
			}
		}
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
