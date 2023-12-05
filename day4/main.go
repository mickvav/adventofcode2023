package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Line struct {
	winning map[int]bool
	hand    []int
	number  int
}
type Copies map[int]int

func (c Copies) win(current, score int) {
	for p:=current + 1; p<=current + score; p++ {
		if _,ok := c[p]; !ok {
			c[p] = c[current]
		} else {
			c[p] += c[current]
		}
	}
}

func (c Copies) take(current int) {
	if _,ok := c[current]; !ok {
		c[current] = 1
	} else {
		c[current] += 1
	}
}

func (c Copies) count() int {
	res := 0
	for _, v := range c {
		res += v
	}
	return res
}
var lineRe = regexp.MustCompile(`Card *([0-9]+): (.*) \| (.*)`)

func ParseLine(s string) (Line, error) {
	l := Line{
		winning: map[int]bool{},
		hand:    []int{},
		number:  0,
	}
	r := lineRe.FindStringSubmatch(s)
	if r == nil {
		return l, fmt.Errorf("parse error for %s", s)
	}
	l.number, _ = strconv.Atoi(r[1])
	ws := strings.Split(r[2], " ")
	for _, wss := range ws {
		if wsi, err := strconv.Atoi(wss); err == nil {
			l.winning[wsi] = true
		}
	}
	hs := strings.Split(r[3], " ")
	for _, hss := range hs {
		if hsi, err := strconv.Atoi(hss); err == nil {
			l.hand = append(l.hand, hsi)
		}
	}
	return l, nil
}

func (l *Line) Score() int {
	s := 0
	for _, v := range l.hand {
		if l.winning[v] {
			if s == 0 {
				s = 1
				continue
			}
			s = s * 2
		}
	}
	return s
}

func (l *Line) Matches() int {
	s := 0
	for _, v := range l.hand {
		if l.winning[v] {
			s+=1
		}
	}
	return s
}

func Readlines(file io.Reader) (int, int) {
	res := 0
	copies := Copies{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		s, e := ParseLine(l)
		if e == nil {
			res += s.Score()
			copies.take(s.number)
			copies.win(s.number, s.Matches())	
		} else {
			fmt.Printf("Problem parsing: %s ", l)
		}
	}
	return res, copies.count()
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Ups")
	}
	defer file.Close()
	f, f1 := Readlines(file)
	fmt.Println(strconv.Itoa(f), strconv.Itoa(f1))
}
