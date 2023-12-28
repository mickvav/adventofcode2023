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
	m      []string
	v      [][]int
	marks  map[state]bool
	values map[state]int
}

type state struct {
	row, col     int
	prevdir      byte
	prevdistance byte
}

func (s state) String() string {
	return fmt.Sprintf("[%d,%d %c %d]", s.row, s.col, s.prevdir, s.prevdistance)
}
func (inp *Input) TraceBack2(slast state) {
	pts := map[dir]byte{}
	s := slast
	for s2, v := range inp.values {
		fmt.Printf("== %s %d\n", s2.String(), v)
	}
	value := inp.values[s]
	for s.row != 0 || s.col != 0 {
		pt := dir{col: s.col, row: s.row}
		value -= inp.v[s.row][s.col]
		pts[pt] = s.prevdir
		d := ltrs[s.prevdir].Reverse()
		if s.prevdistance > 0 {
			s.prevdistance--
			s.col = s.col + d.col
			s.row = s.row + d.row
		} else {
		lettersLoop:
			for k, d2 := range ltrs {
				if k != s.prevdir {
					s1 := state{
						col:     s.col - d2.col,
						row:     s.row - d2.row,
						prevdir: k,
					}
					for prevdistance := 3; prevdistance <= 10; prevdistance++ {
						s1.prevdistance = byte(prevdistance)
						fmt.Println(s1.String())
						if vf, ok := inp.values[s1]; ok {
							if vf == value {
								s = s1
								fmt.Printf("! %v: %d", s1, vf)
								break lettersLoop
							} else {
								fmt.Printf("- %v: %d", s1, vf)
							}
						}
					}
				}
			}
		}
	}
	for i, v := range inp.v {
		for i2, _ := range v {
			if b, ok := pts[dir{row: i, col: i2}]; ok {
				fmt.Printf("%c", b)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:      []string{},
		v:      [][]int{},
		marks:  map[state]bool{},
		values: map[state]int{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		vl := []int{}
		for _, v := range strings.Split(line, "") {
			vv, _ := strconv.Atoi(v)
			vl = append(vl, vv)
		}
		res.v = append(res.v, vl)
	}
	return res
}

type dir struct {
	row, col int
}

var ltrs = map[byte]dir{
	'N': {
		row: -1,
		col: 0,
	},
	'S': {
		row: 1,
		col: 0,
	},
	'E': {
		row: 0,
		col: 1,
	},
	'W': {
		row: 0,
		col: -1,
	},
}

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

func (inp Input) Has(row, col int) bool {
	return row < len(inp.m) && row >= 0 && col < len(inp.m[0]) && col >= 0
}

func (s state) Neighbours(inp Input) []state {
	res := []state{}
	for k, d := range ltrs {
		if k == s.prevdir {
			if s.prevdistance < 3 {
				if inp.Has(d.row+s.row, d.col+s.col) {
					res = append(res, state{
						row:          d.row + s.row,
						col:          d.col + s.col,
						prevdir:      k,
						prevdistance: s.prevdistance + 1,
					})
				}
			}
		} else {
			if d.Reverse() != ltrs[s.prevdir] {
				if inp.Has(d.row+s.row, d.col+s.col) {
					res = append(res, state{
						row:          d.row + s.row,
						col:          d.col + s.col,
						prevdir:      k,
						prevdistance: 1,
					})
				}
			}
		}
	}
	return res
}

func (s state) Neighbours2(inp Input) []state {
	res := []state{}
	for k, d := range ltrs {
		if k == s.prevdir {
			if s.prevdistance <= 9 {
				if inp.Has(d.row+s.row, d.col+s.col) {
					res = append(res, state{
						row:          d.row + s.row,
						col:          d.col + s.col,
						prevdir:      k,
						prevdistance: s.prevdistance + 1,
					})
				}
			}
		} else {
			if d.Reverse() != ltrs[s.prevdir] && s.prevdistance >= 4 {
				if inp.Has(d.row+s.row, d.col+s.col) {
					res = append(res, state{
						row:          d.row + s.row,
						col:          d.col + s.col,
						prevdir:      k,
						prevdistance: 1,
					})
				}
			}
		}
	}
	return res
}

func (inp *Input) fillFirstState() {
	inp.values = map[state]int{
		state{
			row:          0,
			col:          0,
			prevdir:      0,
			prevdistance: 0,
		}: 0,
	}
}
func (inp *Input) fillFirstState2() {
	inp.values = map[state]int{
		{
			row:          0,
			col:          0,
			prevdir:      'E',
			prevdistance: 0,
		}: 0,
		{
			row:          0,
			col:          0,
			prevdir:      'S',
			prevdistance: 0,
		}: 0,
	}
}

func (inp *Input) lastStates() []state {
	res := []state{}
	for k := range ltrs {
		for pd := byte(0); pd <= 3; pd++ {
			res = append(res, state{
				row:          len(inp.m) - 1,
				col:          len(inp.m[0]) - 1,
				prevdir:      k,
				prevdistance: pd,
			})
		}
	}
	return res
}

func (inp *Input) lastStates2() []state {
	res := []state{}
	for k := range ltrs {
		for pd := byte(4); pd <= 10; pd++ {
			res = append(res, state{
				row:          len(inp.m) - 1,
				col:          len(inp.m[0]) - 1,
				prevdir:      k,
				prevdistance: pd,
			})
		}
	}
	return res
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	inp.fillFirstState()
	for {
		ops := 0
		for s, v := range inp.values {
			for _, sn := range s.Neighbours(inp) {
				n1proposed := inp.v[sn.row][sn.col] + v
				if v1, ok := inp.values[sn]; ok {
					if n1proposed < v1 {
						inp.values[sn] = n1proposed
						ops++
					}
				} else {
					inp.values[sn] = n1proposed
					ops++
				}
			}
		}
		if ops == 0 {
			break
		}
	}
	for _, s := range inp.lastStates() {
		if vf, ok := inp.values[s]; ok {
			if s1 == 0 || vf < s1 {
				s1 = vf
			}
		}
	}

	fmt.Println("s1", s1)
	inp.fillFirstState2()
	for {
		ops := 0
		for s, v := range inp.values {
			for _, sn := range s.Neighbours2(inp) {
				n1proposed := inp.v[sn.row][sn.col] + v
				if v1, ok := inp.values[sn]; ok {
					if n1proposed < v1 {
						inp.values[sn] = n1proposed
						ops++
					}
				} else {
					inp.values[sn] = n1proposed
					ops++
				}
			}
		}
		if ops == 0 {
			break
		}
	}
	selectedLastState := state{}
	for _, s := range inp.lastStates2() {
		if vf, ok := inp.values[s]; ok {
			if s2 == 0 || vf < s2 {
				s2 = vf
				selectedLastState = s
			}
		}
	}
	inp.TraceBack2(selectedLastState)

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
