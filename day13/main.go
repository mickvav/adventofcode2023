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
	m [][][]byte
	s []Slide
}

type Slide struct {
	cols []uint32
	rows []uint32
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m: [][][]byte{},
		s: []Slide{},
	}
	sm := make([][]byte, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			res.m = append(res.m, sm)
			res.s = append(res.s, ReadSlide(sm))
			sm = make([][]byte, 0)
			continue
		}
		sm = append(sm, []byte(line))
	}
	res.m = append(res.m, sm)
	res.s = append(res.s, ReadSlide(sm))
	return res
}

func ReadSlide(b [][]byte) Slide {
	s := Slide{
		cols: []uint32{},
		rows: []uint32{},
	}
	for _, v := range b {
		var u uint32 = 0
		for i2, v2 := range v {
			if v2 == '#' {
				u |= 1 << uint32(i2)
			}
		}
		s.rows = append(s.rows, u)
	}
	for i2, _ := range b[0] {
		var u uint32 = 0
		for i, _ := range b {
			if b[i][i2] == '#' {
				u |= 1 << uint32(i)
			}
		}
		s.cols = append(s.cols, u)
	}
	return s
}

func (s Slide) HorizontalReflection() int {
outer:
	for i := 1; i < len(s.rows); i++ {
		j := i - 1
		k := i
		for j >= 0 && k < len(s.rows) {
			if s.rows[j] != s.rows[k] {
				continue outer
			}
			j--
			k++
		}
		return i
	}
	return 0
}

func (s Slide) NewHorizontalReflection(oldrow int) int {
outer:
	for i := 1; i < len(s.rows); i++ {
		if i == oldrow {
			continue
		}
		j := i - 1
		k := i
		for j >= 0 && k < len(s.rows) {
			if s.rows[j] != s.rows[k] {
				continue outer
			}
			j--
			k++
		}
		return i
	}
	return 0
}

func (s Slide) VerticalReflection() int {
outer:
	for i := 1; i < len(s.cols); i++ {
		j := i - 1
		k := i
		for j >= 0 && k < len(s.cols) {
			if s.cols[j] != s.cols[k] {
				continue outer
			}
			j--
			k++
		}
		return i
	}
	return 0
}

func (s Slide) NewVerticalReflection(oldcol int) int {
outer:
	for i := 1; i < len(s.cols); i++ {
		if i == oldcol {
			continue
		}
		j := i - 1
		k := i
		for j >= 0 && k < len(s.cols) {
			if s.cols[j] != s.cols[k] {
				continue outer
			}
			j--
			k++
		}
		return i
	}
	return 0
}

func (s Slide) Smudge(row, col int) {
	s.rows[row] ^= (1 << uint32(col))
	s.cols[col] ^= (1 << uint32(row))
}

func (s Slide) Print() {
	for i, v := range s.rows {
		fmt.Printf("%02d ", i)
		for i2, _ := range s.cols {
			if v&(1<<uint32(i2)) != 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func (s Slide) PrintWithSmudge(row, col int) {
	for i, v := range s.rows {
		fmt.Printf("%02d ", i)
		for i2, _ := range s.cols {
			if v&(1<<uint32(i2)) != 0 {
				if i == row && i2 == col {
					fmt.Print("S")
				} else {
					fmt.Print("#")
				}
			} else {
				if i == row && i2 == col {
					fmt.Print("s")
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Print("\n")
	}
}

func (s Slide) SmudgedReflections(ovs, ohs int) (vs, hs int) {
	for r, _ := range s.rows {
		/*		if ohs > 0 {
				refsize := len(s.rows) - ohs
				if r < ohs-refsize {
					continue
				}
			}*/
		for c, _ := range s.cols {
			/*if ovs > 0 {
				refsize := len(s.cols) - ovs
				if c < ovs-refsize {
					continue
				}
			}*/
			s.Smudge(r, c)
			vs = s.NewVerticalReflection(ovs)
			hs = s.NewHorizontalReflection(ohs)
			if (vs != 0) || hs != 0 {
				s.PrintWithSmudge(r, c)
				s.Smudge(r, c)
				return vs, hs
			}
			s.Smudge(r, c)
		}
	}
	return 0, 0
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	for _, s := range inp.s {
		v := s.VerticalReflection()
		h := s.HorizontalReflection()
		s1 += v + h*100
		v, h = s.SmudgedReflections(v, h)
		s2 += v + h*100
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
