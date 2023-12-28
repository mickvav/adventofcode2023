package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type point struct {
	row, col int
}
type Input struct {
	m         [][]byte
	emptyrows []int
	emptycols []int
	galaxies  []point
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:         [][]byte{},
		emptyrows: []int{},
		emptycols: []int{},
		galaxies:  []point{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, []byte(line))
	}
	return res
}

func (inp *Input) FindEmpty() {
	for i, v := range inp.m {
		empty := true
		for _, v2 := range v {
			if v2 != '.' {
				empty = false
			}
		}
		if empty {
			inp.emptyrows = append(inp.emptyrows, i)
		}
	}
	for i := range inp.m[0] {
		empty := true
		for i2 := range inp.m {
			if inp.m[i2][i] != '.' {
				empty = false
			}
		}
		if empty {
			inp.emptycols = append(inp.emptycols, i)
		}
	}
}

func (inp *Input) InflateGalaxies(factor int) {
	for ip, p := range inp.galaxies {
		er := 0
		for _, v := range inp.emptyrows {
			if v < p.row {
				er++
			} else {
				break
			}
		}
		inp.galaxies[ip].row += er * (factor - 1)
		ec := 0
		for _, v := range inp.emptycols {
			if v < p.col {
				ec++
			} else {
				break
			}
		}
		inp.galaxies[ip].col += ec * (factor - 1)
	}
}
func (inp *Input) InflateY() {
	n := [][]byte{}
	for _, v := range inp.m {
		empty := true
		for _, v2 := range v {
			if v2 != '.' {
				empty = false
			}
		}
		n = append(n, v)
		if empty {
			n = append(n, v)
		}
	}
	inp.m = n
}
func (inp *Input) Transpose() {
	n := [][]byte{}
	for col, _ := range inp.m[0] {
		c := []byte{}
		for _, v2 := range inp.m {
			c = append(c, v2[col])
		}
		n = append(n, c)
	}
	inp.m = n
}

func (inp *Input) Inflate() {
	inp.InflateY()
	inp.Transpose()
	inp.InflateY()
	inp.Transpose()
}

func (inp *Input) FindGalaxies() {
	for row, v := range inp.m {
		for col, v2 := range v {
			if v2 == '#' {
				inp.galaxies = append(inp.galaxies, point{row: row, col: col})
			}
		}
	}
}

func Abs(i int) int {
	if i > 0 {
		return i
	}
	return -i
}
func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	inp.FindEmpty()
	inp.FindGalaxies()
	inp.InflateGalaxies(2)
	for i, p := range inp.galaxies {
		for i2, p2 := range inp.galaxies {
			if i != i2 {
				s1 += Abs(p.col-p2.col) + Abs(p.row-p2.row)
			}
		}
	}
	inp.galaxies = []point{}
	inp.FindGalaxies()
	inp.InflateGalaxies(1000000)
	for i, p := range inp.galaxies {
		for i2, p2 := range inp.galaxies {
			if i != i2 {
				s2 += Abs(p.col-p2.col) + Abs(p.row-p2.row)
			}
		}
	}

	return s1 / 2, s2 / 2
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
