package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

type field struct {
	c []string
}

func (f *field) Process() int {
	res := 0
	NumberRe := regexp.MustCompile("([0-9]+)")
	for row, v := range f.c {
		findings := NumberRe.FindAllStringIndex(v, -1)
		if findings == nil {
			continue
		}
		for _, finding := range findings {
			if f.HasAdjacentChars(row, finding[0], finding[1]) {
				i, _ := strconv.Atoi(v[finding[0]:finding[1]])
				res += i
			}
		}
	}
	return res
}

func (f *field) IsGear(row, col int) (is bool, r, c []int) {
	r = []int{}
	c = []int{}
	for rp := max(row-1, 0); rp < min(row+2, len(f.c)); rp++ {
		fl := false
		for cp := max(col-1, 0); cp < min(col+2, len(f.c[rp])); cp++ {
			v := f.c[rp][cp]
			if v >= '0' && v <= '9' {
				if fl {
					continue
				}
				r = append(r, rp)
				c = append(c, cp)
				fl = true
				continue
			}
			fl = false
		}
	}
	return len(r) == 2, r, c
}

func (f *field) GetPartNumberAt(row, col int) int {
	r := f.c[row]
	var (
		c1 int
		c2 int
	)
	for c1 = col; c1 >= 0 && r[c1] >= '0' && r[c1] <= '9'; c1-- {
	}
	for c2 = col; c2 < len(r) && r[c2] >= '0' && r[c2] <= '9'; c2++ {
	}
	s := r[c1+1 : c2]
	v, e := strconv.Atoi(s)
	if e != nil {
		panic(e)
	}
	return v
}

func (f *field) FindAllGears() int {
	res := 0
	for row, v := range f.c {
		for col := 0; col < len(v); col++ {
			b := v[col]
			if b == '*' {
				is, r, c := f.IsGear(row, col)
				if is {
					pn1 := f.GetPartNumberAt(r[0], c[0])
					pn2 := f.GetPartNumberAt(r[1], c[1])
					gr := pn1 * pn2
					res += gr
				}
			}
		}
	}
	return res
}
func isNotMarker(v byte) bool {
	return (v >= '0' && v <= '9') || (v == '.')
}

func (f *field) HasAdjacentChars(row, start, end int) bool {
	if row > 0 {
		for p := max(start-1, 0); p < min(end+1, len(f.c[row-1])); p++ {
			if isNotMarker(f.c[row-1][p]) {
				continue
			}
			return true
		}
	}
	if row+1 < len(f.c) {
		for p := max(start-1, 0); p < min(end+1, len(f.c[row+1])); p++ {
			if isNotMarker(f.c[row+1][p]) {
				continue
			}
			return true
		}
	}
	if start > 0 && !isNotMarker(f.c[row][start-1]) {
		return true
	}
	if end < len(f.c[row]) && !isNotMarker(f.c[row][end]) {
		return true
	}
	return false
}

func Readlines(file io.Reader) *field {
	f := new(field)
	f.c = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		f.c = append(f.c, s)
	}
	return f
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Ups")
	}
	defer file.Close()
	f := Readlines(file)
	fmt.Println(strconv.Itoa(f.Process()))
	fmt.Println(strconv.Itoa(f.FindAllGears()))
}
