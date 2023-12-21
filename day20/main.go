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
	m          map[string]*Module
	pulseQueue *Queue
	high       int64
	low        int64
}

type Module struct {
	t     byte
	state byte
	links []string
	inps  map[string]byte
}

type pulseSignal struct {
	source, destination string
	value               byte
}
type Queue []pulseSignal

func (self *Queue) Push(x pulseSignal) {
	*self = append(*self, x)
}

func (self *Queue) Pop() pulseSignal {
	h := *self
	var el pulseSignal
	l := len(h)
	el, *self = h[0], h[1:l]
	// Or use this instead for a Stack
	// el, *self = h[l-1], h[0:l-1]
	return el
}

func NewQueue() *Queue {
	return &Queue{}
}

func ReadModule(s string) (string, Module) {
	res := Module{
		t:     0,
		state: 0,
		links: []string{},
		inps:  map[string]byte{},
	}
	p := strings.Split(s, " -> ")
	n := p[0]
	if p[0][0] == '%' || p[0][0] == '&' {
		res.t = p[0][0]
		n = p[0][1:]
	}
	res.links = strings.Split(p[1], ", ")
	return n, res
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:          map[string]*Module{},
		pulseQueue: NewQueue(),
		high:       0,
		low:        0,
	}
	for scanner.Scan() {
		line := scanner.Text()
		s, m := ReadModule(line)
		res.m[s] = &m
	}
	for k, m := range res.m {
		for _, v := range m.links {
			if _, ok := res.m[v]; ok {
				res.m[v].inps[k] = 0
			}
		}
	}
	return res
}

func (inp *Input) Pulse(target string, pulse byte, source string) {
	// fmt.Println(source, " -", pulse, "-> ", target)
	if pulse == 0 {
		inp.low++
	} else {
		inp.high++
	}
	//	m := inp.m[target]

	if m, ok := inp.m[target]; !ok {
		return
	} else {
		switch m.t {
		case 0:
			for _, v := range m.links {
				inp.pulseQueue.Push(pulseSignal{source: target, value: pulse, destination: v})
			}

		case '%':
			if pulse == 0 {
				newstate := 1 - m.state
				m.state = newstate
				for _, v := range m.links {
					inp.pulseQueue.Push(pulseSignal{source: target, value: newstate, destination: v})
				}
			}
		case '&':
			if _, ok := m.inps[source]; ok {
				m.inps[source] = pulse
			}
			out := byte(1)
			for _, v := range m.inps {
				if v == 1 && out == 1 {
					continue
				} else {
					out = 0
					break
				}
			}
			for _, v := range m.links {
				inp.pulseQueue.Push(pulseSignal{source: target, value: 1 - out, destination: v})
			}
		}
	}
}

func (inp Input) Count() (int, int64) {
	s1 := 0
	s2 := int64(0)
	s2parts := map[string]int{}
	for i := 1; i < 1000; i++ {
		inp.pulseQueue.Push(pulseSignal{source: "button", value: 0, destination: "broadcaster"})
		for len(*inp.pulseQueue) > 0 {
			ps := inp.pulseQueue.Pop()
			if ps.destination == "rx" && ps.value == 0 && s2 == 0 {
				s2 = int64(i)
			}
			if ps.value == 0 {
				if ps.destination == "pc" || ps.destination == "nd" || ps.destination == "vd" || ps.destination == "tx" {
					fmt.Printf("%s %d", ps.destination, i)
					s2parts[ps.destination] = i
				}
			}

			inp.Pulse(ps.destination, ps.value, ps.source)
		}
	}
	s1 = int(inp.high) * int(inp.low)
	for i := 1000; s2 == 0; i++ {
		inp.pulseQueue.Push(pulseSignal{source: "button", value: 0, destination: "broadcaster"})
		for len(*inp.pulseQueue) > 0 {
			ps := inp.pulseQueue.Pop()
			if ps.destination == "rx" && ps.value == 0 && s2 == 0 {
				s2 = int64(i)
			}
			if ps.value == 0 {
				if ps.destination == "pc" || ps.destination == "nd" || ps.destination == "vd" || ps.destination == "tx" {
					fmt.Printf("%s %d", ps.destination, i)
					s2parts[ps.destination] = i
				}
			}
			inp.Pulse(ps.destination, ps.value, ps.source)
		}
		if len(s2parts) == 4 {
			s2 = int64(s2parts["pc"] * s2parts["nd"] * s2parts["vd"] * s2parts["tx"])
			fmt.Println(s2, " ", s2parts)
			return s1, s2
		}
		if i%100000 == 0 {
			fmt.Println(i)
		}
	}
	return s1, s2
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
