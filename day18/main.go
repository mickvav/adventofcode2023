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
)

type Input struct {
	m      []Line
	p      map[pos]byte
	p2     map[pos]byte
	min1   pos
	max1   pos
	min2   pos
	max2   pos
	x2     []int
	y2     []int
	x2imap map[int]int // coordinate to index map for x2
	y2imap map[int]int // coordinate to index map for y2
}

type Line struct {
	direction  byte
	length     int
	direction2 byte
	length2    int64
}

type pos struct {
	x, y int
}

func BoundingBox(p pos, min *pos, max *pos) {
	if p.x < min.x {
		min.x = p.x
	}
	if p.x > max.x {
		max.x = p.x
	}
	if p.y < min.y {
		min.y = p.y
	}
	if p.y > max.y {
		max.y = p.y
	}
}
func (p pos) Add(d dir) pos {
	return pos{p.x + d.x, p.y + d.y}
}
func (p pos) AddMultiplied(d dir, a int) pos {
	return pos{p.x + d.x*a, p.y + d.y*a}
}

type dir struct {
	x, y int
}

func (d dir) X() bool {
	return d.x != 0
}

var ltrs = map[byte]dir{
	'U': {x: 0, y: 1},
	'D': {x: 0, y: -1},
	'L': {x: -1, y: 0},
	'R': {x: 1, y: 0},
}
var ltrtobitmapstart = map[byte]byte{
	// 1
	// 0
	'U': 0b01,
	'D': 0b10,
	'L': 0,
	'R': 0,
}
var ltrtobitmapend = map[byte]byte{
	'U': 0b10,
	'D': 0b01,
	'L': 0,
	'R': 0,
}

var ltrtobitmapstart2 = map[byte]byte{
	// 3  U
	// 2  U
	// 1     D
	//
	//     123
	'U': 0b011,
	'D': 0b100,
	'L': 0,
	'R': 0,
}
var ltrtobitmapend2 = map[byte]byte{
	// 3     D
	// 2     D
	// 1   U
	//
	//     123
	'U': 0b100,
	'D': 0b011,
	'L': 0,
	'R': 0,
}

var htoltr = map[byte]byte{
	'0': 'R',
	'1': 'D',
	'2': 'L',
	'3': 'U',
}

func ReadLine(s string) Line {
	res := Line{
		direction:  0,
		length:     0,
		direction2: 0,
		length2:    0,
	}
	p := strings.Split(s, " ")
	res.direction = p[0][0]
	res.length, _ = strconv.Atoi(p[1])
	res.length2, _ = strconv.ParseInt(p[2][2:7], 16, 32)
	res.direction2 = htoltr[p[2][7]]
	return res
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:      []Line{},
		p:      map[pos]byte{},
		p2:     map[pos]byte{},
		min1:   pos{},
		max1:   pos{},
		min2:   pos{},
		max2:   pos{},
		x2:     []int{},
		y2:     []int{},
		x2imap: map[int]int{},
		y2imap: map[int]int{},
	}
	ps := pos{}
	ps2 := pos{}
	res.p[ps] = 0
	res.p2[ps2] = 0
	x2m := map[int]bool{0: true}
	y2m := map[int]bool{0: true}
	for scanner.Scan() {
		line := scanner.Text()
		parsedline := ReadLine(line)
		res.m = append(res.m, parsedline)
		d := ltrs[parsedline.direction]
		for i := 0; i < parsedline.length; i++ {
			res.p[ps] |= ltrtobitmapstart[parsedline.direction]
			ps = ps.Add(d)
			res.p[ps] |= ltrtobitmapend[parsedline.direction]
		}
		BoundingBox(ps, &res.min1, &res.max1)
		d = ltrs[parsedline.direction2]
		ps2 = ps2.AddMultiplied(d, int(parsedline.length2))
		x2m[ps2.x] = true
		y2m[ps2.y] = true
		BoundingBox(ps2, &res.min2, &res.max2)
	}
	ps2 = pos{0, 0}
	res.x2 = make([]int, 0, len(x2m))
	for k := range x2m {
		res.x2 = append(res.x2, k)
	}
	sort.Slice(res.x2, func(i, j int) bool {
		return res.x2[i] < res.x2[j]
	})
	for i, v := range res.x2 {
		res.x2imap[v] = i
	}
	res.y2 = make([]int, 0, len(y2m))
	for k := range y2m {
		res.y2 = append(res.y2, k)
	}
	sort.Slice(res.y2, func(i, j int) bool {
		return res.y2[i] < res.y2[j]
	})
	for i, v := range res.y2 {
		res.y2imap[v] = i
	}

	for _, l := range res.m {
		d := ltrs[l.direction2]
		ps2next := ps2.AddMultiplied(d, int(l.length2))
		pi := ps2
		// (ps2next)
		// |
		// |
		// (ps2)
		res.p2[pi] |= ltrtobitmapstart2[l.direction2]
		for xidx := res.x2imap[ps2.x] + d.x; xidx != res.x2imap[ps2next.x]; xidx += d.x {
			res.p2[pi] |= ltrtobitmapstart2[l.direction2]
			pi = pos{x: res.x2[xidx], y: ps2.y}
			res.p2[pi] |= ltrtobitmapend2[l.direction2]
		}
		pi = ps2
		for yidx := res.y2imap[ps2.y] + d.y; yidx != res.y2imap[ps2next.y]; yidx += d.y {
			res.p2[pi] |= ltrtobitmapstart2[l.direction2]
			pi = pos{x: ps2.x, y: res.y2[yidx]}
			res.p2[pi] |= 0b111
		}
		res.p2[ps2next] |= ltrtobitmapend2[l.direction2]
		ps2 = ps2next
	}
	return res
}

func (inp Input) Count() (int, int64) {
	s1 := 0
	s2 := int64(0)
	for y := inp.min1.y; y <= inp.max1.y; y++ {
		in := byte(0)
		for x := inp.min1.x; x <= inp.max1.x; x++ {
			if v, ok := inp.p[pos{x: x, y: y}]; ok {
				s1++
				fmt.Print("#")
				in ^= v
			} else {
				if in == 0b11 {
					fmt.Print("O")
					s1++
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
	for y2idx, y2 := range inp.y2 {
		in := byte(0)
		lengthy := 1
		if y2idx < len(inp.y2)-1 {
			lengthy = inp.y2[y2idx+1] - y2
		}
		for x2idx, x2 := range inp.x2 {
			lengthx := 1
			if x2idx < len(inp.x2)-1 {
				lengthx = inp.x2[x2idx+1] - x2
			}
			oldin := in
			if v, ok := inp.p2[pos{x: x2, y: y2}]; ok {
				in ^= v
			}
			switch in {
			case 0b100:
				s2 += int64(lengthx)
				if oldin == 0b111 {
					fmt.Print("L")
					s2 += int64(lengthy) - 1
				} else {
					fmt.Print("_")
				}
			case 0b111:
				s2 += int64(lengthx) * int64(lengthy)
				fmt.Print("#")
			case 0b011:
				s2 += int64(lengthx) * int64(lengthy)
				fmt.Print("±")
			default:
				switch oldin {
				case 0b111:
					s2 += int64(lengthy)
					fmt.Print("|")
				case 0b100:
					s2 += 1
					fmt.Print(",")
				case 0b011:
					s2 += int64(lengthy)
					fmt.Print("/")
				default:
					fmt.Print(".")
				}
			}

		}
		fmt.Println()
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
