package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Input struct {
	m  []Line
	m5 []Line
	mh mMatcher
}

type Line struct {
	orig      string
	length    uint64
	masks     []mask
	checksums []int
}

type mask struct {
	unknown     uint16
	damaged     uint16
	operational uint16
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

func requiredMaskLength(l int) int {
	if l%16 == 0 {
		return l / 16
	}
	return 1 + l/16
}
func ReadLine(s string) Line {
	res := Line{
		orig:      "",
		length:    0,
		masks:     []mask{},
		checksums: []int{},
	}
	p := strings.Split(s, " ")
	res.orig = p[0]
	for _, v := range strings.Split(p[1], ",") {
		cs, _ := strconv.Atoi(v)
		res.checksums = append(res.checksums, cs)
	}
	res.length = uint64(len(p[0]))
	res.masks = make([]mask, requiredMaskLength(len(p[0])))
	for i, v := range strings.Split(p[0], "") {
		m := i / 16
		bp := i % 16
		switch v {
		case "?":
			res.masks[m].unknown |= 1 << bp
		case "#":
			res.masks[m].damaged |= 1 << bp
		case ".":
			res.masks[m].operational |= 1 << bp
		}
	}
	if len(p[0])%16 != 0 {
		m := len(p[0]) / 16
		for bp := len(p[0]) % 16; bp < 16; bp++ {
			res.masks[m].operational |= 1 << bp
		}

	}
	return res
}

type guess uint16

type mMatcher struct {
	h map[mask]map[intStructure]uint16
}

func (is intStructure) String() string {
	res := "["
	for i := 0; i < len(is); i++ {
		res += strconv.Itoa(int(is[i])) + " "
	}
	res += "]"
	return res
}
func (m *mMatcher) Mask(msk mask) map[intStructure]uint16 {
	if r, ok := m.h[msk]; ok {
		return r
	}
	m.h[msk] = make(map[intStructure]uint16)
	checksum := 0
	for i := 0; i < 256*256; i++ {
		if msk.Match(guess(i)) {
			m.h[msk][IntervalStructure[i]]++
			checksum++
		}
	}
	fmt.Printf("Mask : %v, checksum : %d\n", msk, checksum)
	return m.h[msk]
}

type intStructure string

func (is intStructure) Concat(b intStructure) intStructure {
	if is[len(is)-1] != 0 {
		if b[0] != 0 {
			return is[0:len(is)-1] + intStructure([]byte{b[0] + is[len(is)-1]}) + b[1:]
		} else {
			return is + b[1:]
		}
	} else {
		if b[0] != 0 {
			return is[:len(is)-1] + b

		} else {
			if len(is) == 1 {
				return b
			}
			return is[:len(is)-1] + b[1:]
		}
	}
}

func (is intStructure) IsStartOf(m []int) bool {
	j := 0
	for i := 0; i < len(is); i++ {
		if is[i] != 0 {
			if j >= len(m) {
				return false
			}
			if is[i] != byte(m[j]) {
				if is[i] > byte(m[j]) {
					return false
				}
				if i != len(is)-1 {
					return false
				}
			}
			j++
		}
	}
	return true
}

func (is intStructure) Matches(m []int) bool {
	j := 0
	for i := 0; i < len(is); i++ {
		if is[i] != 0 {
			if j >= len(m) {
				return false
			}
			if is[i] != byte(m[j]) {
				return false
			}
			j++
		}
	}
	return j == len(m)
}

var IntervalStructure = [256 * 256]intStructure{}

func (g guess) calcIntervalStructure() intStructure {
	state := 1
	intlen := byte(0)
	res := []byte{}
	for i := uint16(0); i < 16; i++ {
		v := 1 & (uint16(g) >> i)
		if v == 1 {
			switch state {
			case 1:
				intlen++
				continue
			case 0:
				intlen = 1
				state = 1
			}
		} else {
			if state == 1 {
				res = append(res, intlen)
				intlen = 0
				state = 0
			}
		}
	}
	res = append(res, intlen)
	return intStructure(res)
}

func (m *mask) Match(damaged guess) bool {
	if m.operational&uint16(damaged) > 0 {
		return false
	}
	if m.damaged&^uint16(damaged) > 0 {
		return false
	}
	return true
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:  []Line{},
		m5: []Line{},
		mh: mMatcher{
			h: map[mask]map[intStructure]uint16{},
		},
	}
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, ReadLine(line))
		res.m5 = append(res.m5, ReadLine5x(line))
	}
	return res
}

func countmhs(inp intStructure, mhs []map[intStructure]uint16, l *Line) uint64 {
	if len(mhs) == 1 {
		s := uint64(0)
		for is, v := range mhs[0] {
			if inp.Concat(is).Matches(l.checksums) {
				s += uint64(v)
			}
		}
		return s
	} else {
		s := uint64(0)
		hcs := map[intStructure]uint64{}
		for is, v := range mhs[0] {
			cs := inp.Concat(is)
			if cmhs, ok := hcs[cs]; ok {
				s += uint64(v) * cmhs
			} else {
				if cs.IsStartOf(l.checksums) {
					hcs[cs] = countmhs(cs, mhs[1:], l)
					s += uint64(v) * hcs[cs]
				} else {
					hcs[cs] = 0
				}
			}
		}
		return s
	}
}
func (inp Input) Count() (int, uint64) {
	s1 := 0
	s2 := uint64(0)
	for _, l := range inp.m {
		mhs := []map[intStructure]uint16{}
		for _, m := range l.masks {
			mh := inp.mh.Mask(m)
			mhs = append(mhs, mh)
		}
		fmt.Print(".")
		s1 += int(countmhs(intStructure([]byte{0}), mhs, &l))
	}
	for idx, l := range inp.m5 {
		mhs := []map[intStructure]uint16{}
		n := time.Now()
		for _, m := range l.masks {
			mh := inp.mh.Mask(m)
			mhs = append(mhs, mh)
		}
		s2 += countmhs(intStructure([]byte{0}), mhs, &l)
		fmt.Println(time.Since(n), " ", idx, " ", l.length)
	}
	return s1, s2
}

func Readlines(file io.Reader) (int, uint64) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	return inp.Count()
}

func init() {

	for i := 0; i < 256*256; i++ {
		IntervalStructure[i] = guess(i).calcIntervalStructure()
	}
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
