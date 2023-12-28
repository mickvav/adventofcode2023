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
	m  []Line
	m5 []Line
}

type uint128 struct {
	hi, lo uint64
}
type Line struct {
	orig        string
	length      uint64
	unknowns    uint128
	damaged     uint128
	operational uint128
	checksums   []uint64
}

type LineIterator struct {
	length        uint64
	damagedlength uint64
	checksums     []uint64
	state         []uint64
}

func (l Line) StartIteration() LineIterator {
	res := LineIterator{
		length:        l.length,
		damagedlength: 0,
		checksums:     l.checksums,
		state:         make([]uint64, len(l.checksums)+1),
	}
	for i := 1; i < len(res.state)-1; i++ {
		res.state[i] = 1
	}
	for _, v := range l.checksums {
		res.damagedlength += v
	}
	res.state[0] = l.length - res.damagedlength - uint64(len(res.state)-2)
	return res
}

func (u *uint128) put1(j uint64) {
	if j < 64 {
		u.lo |= 1 << j
	} else {
		u.hi |= 1 << (j - 64)
	}
}

func (u *uint128) and(u1 *uint128) uint128 {
	return uint128{
		hi: u.hi & u1.hi,
		lo: u.lo & u1.lo,
	}
}

func (u *uint128) andnot(u1 *uint128) uint128 {
	return uint128{
		hi: u.hi & (^u1.hi),
		lo: u.lo & (^u1.lo),
	}
}

func (u *uint128) not0() bool {
	return u.hi != 0 || u.lo != 0
}
func (u *uint128) has1above(j uint64) bool {
	if j < 64 {
		return (u.hi != 0) || ((u.lo >> j) != 0)
	} else {
		return (u.hi >> (j - 64)) != 0
	}
}
func (it *LineIterator) Repr() uint128 {
	res := uint128{}
	j := uint64(0)
	for i, v := range it.state {
		j += v
		if i < len(it.checksums) {
			for k := uint64(0); k < it.checksums[i]; k++ {
				res.put1(j)
				j++
			}
		}
	}
	return res
}

func (it *LineIterator) Step(hint uint128) bool {
	p := 1
	for p = 1; p < len(it.state); p++ {
		it.state[p]++
		ps := uint64(0)
		for p1 := p; p1 < len(it.state); p1++ {
			ps += it.state[p1]
		}
		if it.length < it.damagedlength+ps+uint64(p-1) {
			it.state[p] = 1
			continue
		}
		bs := it.state[p] - 1
		for p1 := p + 1; p1 < len(it.state); p1++ {
			bs += it.state[p1]
			bs += it.checksums[p1-1]
		}
		if hint.has1above(it.length - bs) {
			it.state[p] = 1
			continue
		}
		break
	}
	if p >= len(it.state) {
		return false
	}
	ps := uint64(0)
	for p1 := 1; p1 < len(it.state); p1++ {
		ps += it.state[p1]
	}
	if it.length < it.damagedlength+ps {
		return false
	}
	it.state[0] = it.length - it.damagedlength - ps
	return true
}

func ReadLine5x(s string) Line {
	p := strings.Split(s, " ")
	newpattern := ""
	newlimits := ""
	for i := 0; i < 5; i++ {
		newpattern += p[0] + "?"
		newlimits += p[1] + ","
	}
	return ReadLine(newpattern[0:len(newpattern)-1] + " " + strings.Trim(newlimits, ","))
}

func ReadLine(s string) Line {
	res := Line{
		orig:   "",
		length: 0,
		unknowns: uint128{
			hi: 0,
			lo: 0,
		},
		damaged: uint128{
			hi: 0,
			lo: 0,
		},
		operational: uint128{
			hi: 0,
			lo: 0,
		},
		checksums: []uint64{},
	}
	p := strings.Split(s, " ")
	res.orig = p[0]
	for _, v := range strings.Split(p[1], ",") {
		cs, _ := strconv.Atoi(v)
		res.checksums = append(res.checksums, uint64(cs))
	}
	res.length = uint64(len(p[0]))
	for i, v := range strings.Split(p[0], "") {
		switch v {
		case "?":
			res.unknowns.put1(uint64(i))
		case "#":
			res.damaged.put1(uint64(i))
		case ".":
			res.operational.put1(uint64(i))
		}
	}
	return res
}

func (l Line) Match(damaged uint128) (bool, uint128) {
	v := l.operational.and(&damaged)
	if v.not0() {
		return false, v
	}
	v = l.damaged.andnot(&damaged)
	if v.not0() {
		return false, v
	}
	return true, uint128{hi: 0, lo: 0}
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:  []Line{},
		m5: []Line{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, ReadLine(line))
		res.m5 = append(res.m5, ReadLine5x(line))
	}
	return res
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	for _, l := range inp.m {
		it := l.StartIteration()
		rs := 0
		hint := uint128{}
		ok := true
		for {
			if ok, hint = l.Match(it.Repr()); ok {
				rs++
			}
			if !it.Step(hint) {
				break
			}
		}
		fmt.Println(l.orig, rs)
		s1 += rs
	}
	numJobs := len(inp.m5)
	jobs := make(chan int, numJobs)
	results := make(chan uint64, numJobs)
	for w := 1; w <= 16; w++ {
		go worker(w, &inp, jobs, results)
	}
	for i, _ := range inp.m5 {
		jobs <- i
	}
	close(jobs)
	for i := range inp.m5 {
		rs := <-results
		fmt.Println(i)
		s2 += int(rs)
	}
	/*	for _, l := range inp.m5 {
		it := l.StartIteration()
		rs := 0
		hint := uint128{}
		ok := true
		for {
			if ok, hint = l.Match(it.Repr()); ok {
				rs++
			}
			if !it.Step(hint) {
				break
			}
		}
		fmt.Println(l.orig, rs)
		s2 += rs
	}*/

	return s1, s2
}

func worker(id int, inp *Input, jobs <-chan int, results chan<- uint64) {
	for j := range jobs {
		l := inp.m5[j]
		it := l.StartIteration()
		rs := 0
		hint := uint128{}
		ok := true
		for {
			if ok, hint = l.Match(it.Repr()); ok {
				rs++
			}
			if !it.Step(hint) {
				break
			}
		}
		fmt.Println(id, j, l.orig, rs)
		results <- uint64(rs)
	}
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
