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
	s     string
	steps []string
}

type Lens struct {
	f int
	l string
}

type Box struct {
	l []Lens
}
type Boxes struct {
	b []*Box
}

func NewBoxes() Boxes {
	res := Boxes{
		b: []*Box{},
	}
	for i := 0; i < 256; i++ {
		res.b = append(res.b, new(Box))
		res.b[i].l = make([]Lens, 0)
	}
	return res
}

func (b *Boxes) FocalPower() int {
	res := 0
	for i, b2 := range b.b {
		for i2, l := range b2.l {
			res += (1 + i) * (1 + i2) * l.f
		}
	}
	return res
}
func (b *Boxes) Print() {
	for i, b2 := range b.b {
		if len(b2.l) > 0 {
			fmt.Print("Bix ", i, ": ")
			for _, l := range b2.l {
				fmt.Print("[" + l.l + " " + strconv.Itoa(l.f) + "]")
			}
			fmt.Print("\n")
		}
	}
}

func (b *Box) dash(l string) {
	pos := -1
	for i, l2 := range b.l {
		if l2.l == l {
			pos = i
			break
		}
	}
	if pos >= 0 {
		b.l = append(b.l[:pos], b.l[pos+1:]...)
	}
}

func (b *Box) eq(l string, f int) {
	for i, l2 := range b.l {
		if l2.l == l {
			b.l[i].f = f
			return
		}
	}
	b.l = append(b.l, Lens{
		f: f,
		l: l,
	})
}

func HolidayHash(s string) byte {
	var c int = 0
	for _, v := range []byte(s) {
		c += int(v)
		c *= 17
		c %= 256
	}
	return byte(c)
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		s:     "",
		steps: []string{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		res.s = line
		res.steps = strings.Split(line, ",")
	}
	return res
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	bx := NewBoxes()
	for _, v := range inp.steps {
		s1 += int(HolidayHash(v))
		p1 := strings.Split(v, "=")
		if len(p1) == 2 {
			f, _ := strconv.Atoi(p1[1])
			bn := HolidayHash(p1[0])
			bx.b[bn].eq(p1[0], f)
		}
		p1 = strings.Split(v, "-")
		if len(p1) == 2 {
			bn := HolidayHash(p1[0])
			bx.b[bn].dash(p1[0])
		}
	}
	s2 = bx.FocalPower()

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
