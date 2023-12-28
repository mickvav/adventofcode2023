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
	m    []string
	b    []*brick
	zmap map[int]map[int]*brick
	// .     z .     rn.
}

type brick struct {
	rn  int
	img [4]uint64
	z   int
	h   int
}
type image uint16

func (b *brick) ximage() image {
	res := uint64(0)
	for y := 0; y < 4; y++ {
		f := 0b1111111111111111 & (b.img[0] >> uint64(16*y))
		f |= 0b1111111111111111 & (b.img[1] >> uint64(16*y))
		f |= 0b1111111111111111 & (b.img[2] >> uint64(16*y))
		f |= 0b1111111111111111 & (b.img[3] >> uint64(16*y))
		res = res | f
	}
	return image(res)
}

// 0b0001
//
//	0001
//	0000
func (b *brick) yimage() image {
	res := uint64(0)
	for y := 0; y < 4; y++ {
		b1 := (0b1111111111111111) & (b.img[0] >> uint64(16*y))
		b2 := (0b1111111111111111) & (b.img[1] >> uint64(16*y))
		b3 := (0b1111111111111111) & (b.img[2] >> uint64(16*y))
		b4 := (0b1111111111111111) & (b.img[3] >> uint64(16*y))
		if b1 > 0 {
			res |= 1 << y
		}
		if b2 > 0 {
			res |= 1 << (4 + y)
		}
		if b3 > 0 {
			res |= 1 << (8 + y)
		}
		if b4 > 0 {
			res |= 1 << (12 + y)
		}
	}
	return image(res)
}

func ReadBrick(s string, rn int) brick {
	p := strings.Split(s, "~")
	p1 := strings.Split(p[0], ",")
	p2 := strings.Split(p[1], ",")
	b := brick{
		rn:  rn,
		img: [4]uint64{0, 0, 0, 0},
		z:   0,
		h:   0,
	}
	p1x, _ := strconv.Atoi(p1[0])
	p1y, _ := strconv.Atoi(p1[1])
	p1z, _ := strconv.Atoi(p1[2])
	p2x, _ := strconv.Atoi(p2[0])
	p2y, _ := strconv.Atoi(p2[1])
	p2z, _ := strconv.Atoi(p2[2])
	b.z = min(p1z, p2z)
	b.h = max(p2z, p1z) - min(p2z, p1z) + 1
	if b.h < 0 {
		b.h = -b.h
	}
	for x := p1x; x <= p2x; x++ {
		for y := p1y; y <= p2y; y++ {
			pos := x + y*16
			switch y {
			case 0, 1, 2, 3:
				b.img[0] |= (1 << pos)
			case 4, 5, 6, 7:
				b.img[1] |= (1 << (pos - 64))
			case 8, 9, 10, 11:
				b.img[2] |= (1 << (pos - 128))
			case 12, 13, 14, 15:
				b.img[3] |= (1 << (pos - 128 - 64))
			}
		}
	}
	return b
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:    []string{},
		b:    []*brick{},
		zmap: map[int]map[int]*brick{},
	}
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		b := ReadBrick(line, row)
		res.b = append(res.b, &b)
		for z := b.z; z < b.z+b.h; z++ {
			if _, ok := res.zmap[z]; ok {
				res.zmap[z][row] = &b
			} else {
				res.zmap[z] = map[int]*brick{row: &b}
			}
		}
		row++
	}
	return res
}

func (b image) Repr() string {
	res := ""
	for i := 0; i < 16; i++ {
		if 0b1&(b>>i) > 0 {
			res += "#"
		} else {
			res += "."
		}
	}
	return res
}

func (inp *Input) Print() {
	keys := make([]int, 0, len(inp.zmap))
	for k := range inp.zmap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, z := range keys {
		ximage := image(0)
		yimage := image(0)
		s := " "
		for k, b := range inp.zmap[z] {
			ximage |= b.ximage()
			yimage |= b.yimage()
			s += strconv.Itoa(k) + " "
		}
		fmt.Printf("%s   %s  (%d: %s)\n", ximage.Repr(), yimage.Repr(), z, s)
	}
	fmt.Println()
}

func modNeg(v, m int) int {
	return (v%m + m) % m
}

func (inp *Input) Zbelow(b *brick) int {
	if b.z == 1 {
		return 0
	}
	for z := b.z - 1; z >= 1; z-- {
		if m, ok := inp.zmap[z]; ok {
			for _, b2 := range m {
				if (b.img[0]&b2.img[0] > 0) || (b.img[1]&b2.img[1] > 0) || (b.img[2]&b2.img[2] > 0) || (b.img[3]&b2.img[3] > 0) {
					dh := b.z - (z + 1)
					if b.z-dh < 1 {
						fmt.Println("!!!!")
					}
					return dh
				}
			}
		}
	}
	return b.z - 1
}

func (inp *Input) ZbelowBut(b *brick, removedbrick *brick) int {
	if b.z == 1 {
		return 0
	}
	for z := b.z - 1; z >= 1; z-- {
		if m, ok := inp.zmap[z]; ok {
			for _, b2 := range m {
				if b2 == removedbrick {
					continue
				}
				if (b.img[0]&b2.img[0] > 0) || (b.img[1]&b2.img[1] > 0) || (b.img[2]&b2.img[2] > 0) || (b.img[3]&b2.img[3] > 0) {
					return b.z - (z + 1)
				}
			}
		}
	}
	return b.z - 1
}

func (inp *Input) Drop(b *brick, h int) {
	for z := b.z; z < b.z+b.h; z++ {
		delete(inp.zmap[z], b.rn)
	}
	b.z = b.z - h
	for z := b.z; z < b.z+b.h; z++ {
		if _, ok := inp.zmap[z]; ok {
			inp.zmap[z][b.rn] = b
		} else {
			inp.zmap[z] = map[int]*brick{b.rn: b}
		}
	}
}

func (inp *Input) DropAll() {
	ops := 1
	for ops > 0 {
		ops = 0
		for _, b := range inp.b {
			zb := inp.Zbelow(b)
			if zb > 0 {
				ops += 1
				inp.Drop(b, zb)
			}
		}
	}
}

func (inp *Input) DropAllUniqueFalls() int {
	ops := 1
	fallingblocks := map[int]bool{}
	for ops > 0 {
		ops = 0
		for _, b := range inp.b {
			zb := inp.Zbelow(b)
			if zb > 0 {
				ops += 1
				inp.Drop(b, zb)
				fallingblocks[b.rn] = true
			}
		}
	}
	return len(fallingblocks)
}

func (inp *Input) WillFallIfDisintegrated(b *brick) int {
	inpNew := Input{
		m:    []string{},
		b:    []*brick{},
		zmap: map[int]map[int]*brick{},
	}
	for _, b2 := range inp.b {
		if b2 != b {
			inpNew.b = append(inpNew.b, &brick{
				rn:  b2.rn,
				img: b2.img,
				z:   b2.z,
				h:   b2.h,
			})
		}
	}
	for _, b2 := range inpNew.b {
		for z := b2.z; z < b2.z+b2.h; z++ {
			if _, ok := inpNew.zmap[z]; ok {
				inpNew.zmap[z][b2.rn] = b2
			} else {
				inpNew.zmap[z] = map[int]*brick{b2.rn: b2}
			}
		}
	}
	return inpNew.DropAllUniqueFalls()
}

func (inp *Input) CanBeDisintegrated(b *brick) bool {
	nextz := b.z + b.h
	if m, ok := inp.zmap[nextz]; !ok {
		return true
	} else {
		for _, b2 := range m {
			if inp.ZbelowBut(b2, b) > 0 {
				return false
			}
		}
		return true
	}
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	inp.DropAll()
	for _, b := range inp.b {
		if inp.CanBeDisintegrated(b) {
			s1++
		} else {
			s2 += inp.WillFallIfDisintegrated(b)
		}
	}
	return s1, s2
}

func Readlines(file io.Reader) (int, int) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	s1, s2 := inp.Count()
	inp.Print()
	return s1, s2
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
