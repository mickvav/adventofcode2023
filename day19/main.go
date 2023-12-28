package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Input struct {
	workflows  map[string][]step
	parts      []part
	attrmap    map[byte]map[int]bool
	attrs      map[byte][]int
	attridxmap map[byte]map[int]int
}

type step struct {
	cond   byte
	value  int
	attr   byte
	target string
}

type part struct {
	x, m, a, s int
}

func (p part) attr(c byte) int {
	switch c {
	case 'x':
		return p.x
	case 'm':
		return p.m
	case 'a':
		return p.a
	case 's':
		return p.s
	}
	return 0
}

func (inp *Input) ProcessMaps() {
	for k, v := range inp.attrmap {
		keys := make([]int, 0, len(v))
		for k2 := range v {
			keys = append(keys, k2)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		inp.attrs[k] = keys
		for i, v2 := range keys {
			inp.attridxmap[k][v2] = i
		}
	}
}

func (inp *Input) ProcessPartRec(p part, wf string) string {
	steps := inp.workflows[wf]
	for _, s2 := range steps {
		if s2.cond == 0 {
			if s2.target == "A" || s2.target == "R" {
				return s2.target
			}
			return inp.ProcessPart(p, s2.target)
		}
		m := false
		v := p.attr(s2.attr)
		switch s2.cond {
		case '>':
			m = v > s2.value
		case '<':
			m = v < s2.value
		}
		if m {
			if s2.target == "A" || s2.target == "R" {
				return s2.target
			}
			return inp.ProcessPart(p, s2.target)
		}
	}
	return ""
}

func (inp *Input) ProcessPart(p part, wf string) string {
	rv := wf
outerloop:
	for rv != "A" && rv != "R" {
		steps := inp.workflows[rv]
		for _, s2 := range steps {
			if s2.cond == 0 {
				rv = s2.target
				continue outerloop
			}
			m := false
			v := p.attr(s2.attr)
			switch s2.cond {
			case '>':
				m = v > s2.value
			case '<':
				m = v < s2.value
			}
			if m {
				rv = s2.target
				continue outerloop
			}
		}
	}
	return rv
}

func ReadWorkflow(s string) (string, []step) {
	res := ""
	st := []step{}

	p := strings.Split(s, "{")
	res = p[0]
	for _, v := range strings.Split(p[1][:len(p[1])-1], ",") {
		stt := strings.Split(v, ":")
		if len(stt) > 1 {
			value, _ := strconv.Atoi(stt[0][2:])
			st = append(st, step{
				cond:   stt[0][1],
				value:  value,
				attr:   stt[0][0],
				target: stt[1],
			})
		} else {
			st = append(st, step{
				cond:   0,
				value:  0,
				attr:   0,
				target: stt[0],
			})
		}
	}
	return res, st
}

func ReadPart(s string) part {
	res := part{}
	for _, v := range strings.Split(s[1:len(s)-1], ",") {
		p := strings.Split(v, "=")
		switch p[0] {
		case "x":
			res.x, _ = strconv.Atoi(p[1])
		case "m":
			res.m, _ = strconv.Atoi(p[1])
		case "a":
			res.a, _ = strconv.Atoi(p[1])
		case "s":
			res.s, _ = strconv.Atoi(p[1])
		}
	}
	return res
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		workflows:  map[string][]step{},
		parts:      []part{},
		attrmap:    map[byte]map[int]bool{'x': {1: true, 4000: true}, 'm': {1: true, 4000: true}, 'a': {1: true, 4000: true}, 's': {1: true, 4000: true}},
		attrs:      map[byte][]int{},
		attridxmap: map[byte]map[int]int{'x': {}, 'm': {}, 'a': {}, 's': {}},
	}
	state := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			state = 1
			continue
		}
		if state == 0 {
			name, steps := ReadWorkflow(line)
			res.workflows[name] = steps
			for _, s := range steps {
				if s.cond != 0 {
					res.attrmap[s.attr][s.value] = true
					res.attrmap[s.attr][s.value+1] = true
					res.attrmap[s.attr][s.value-1] = true
				}
			}
		} else {
			res.parts = append(res.parts, ReadPart(line))
		}
	}
	res.ProcessMaps()
	return res
}

func (inp Input) Count() (int, int64) {
	s1 := 0
	s2 := int64(0)
	for _, p := range inp.parts {
		res := inp.ProcessPart(p, "in")
		if res == "A" {
			s1 += p.s + p.m + p.x + p.a
		}
	}
	numjobs := len(inp.attrs['x']) * len(inp.attrs['m'])
	jobs := make(chan jobdescription, numjobs)
	results := make(chan int64, numjobs)
	for w := 1; w <= 7; w++ {
		go worker(w, &inp, jobs, results)
	}
	for ix, x := range inp.attrs['x'] {
		xlength := 1
		if ix+1 < len(inp.attrs['x']) {
			xlength = inp.attrs['x'][ix+1] - x
		}
		for im, m := range inp.attrs['m'] {
			jobs <- jobdescription{im: im, m: m, x: x, xlength: xlength}
			//			s2 += newFunction(im, &inp, m, x, xlength)
		}
	}
	close(jobs)
	for ix := range inp.attrs['x'] {
		st := time.Now()
		for range inp.attrs['m'] {
			rs := <-results
			s2 += rs
		}
		fmt.Printf("[%d / %d] %s\n", ix, len(inp.attrs['x']), time.Since(st))
	}
	return s1, s2
}

type jobdescription struct {
	im      int
	m       int
	x       int
	xlength int
}

func worker(id int, inp *Input, jobs <-chan jobdescription, results chan<- int64) {
	for j := range jobs {
		results <- newFunction(j.im, inp, j.m, j.x, j.xlength)
	}
}

func newFunction(im int, inp *Input, m int, x int, xlength int) int64 {
	mlength := 1
	s2 := int64(0)
	if im+1 < len(inp.attrs['m']) {
		mlength = inp.attrs['m'][im+1] - m
	}
	marea := mlength * xlength
	for ia, a := range inp.attrs['a'] {
		alength := 1
		if ia+1 < len(inp.attrs['a']) {
			alength = inp.attrs['a'][ia+1] - a
		}
		aarea := marea * alength
		for is, s := range inp.attrs['s'] {
			slength := 1
			if is+1 < len(inp.attrs['s']) {
				slength = inp.attrs['s'][is+1] - s
			}
			res := inp.ProcessPart(part{x: x, m: m, a: a, s: s}, "in")
			if res == "A" {
				s2 += int64(aarea * slength)
			}
		}
	}
	return s2
}

func Readlines(file io.Reader) (int, int64) {
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
