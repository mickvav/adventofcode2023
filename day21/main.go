package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Input struct {
	m     []string
	v     [][]bool
	min   pos
	max   pos
	marks map[pos]bool
}

type pos struct {
	row, col int
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:     []string{},
		v:     [][]bool{},
		marks: map[pos]bool{},
	}
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		vl := []bool{}
		for col, v := range strings.Split(line, "") {
			if v == "." || v == "S" {
				vl = append(vl, true)
			} else {
				vl = append(vl, false)
			}
			if v == "S" {
				res.marks[pos{
					row: row,
					col: col,
				}] = true
				res.min = pos{row: row, col: col}
				res.max = pos{row: row, col: col}
			}
		}
		res.v = append(res.v, vl)
		row++
	}
	return res
}

type bitmap struct {
	buff []byte
	bits int
}

func (b *bitmap) String() string {
	return string(b.buff)
}

func bitMapFromString(s string) bitmap {
	res := bitmap{buff: []byte(s), bits: 0}
	for _, v := range res.buff {
		res.bits += bits.OnesCount8(v)
	}
	return res
}
func (inp *Input) BitMap(offsetrow, offsetcol int) *bitmap {
	res := make([]byte, len(inp.v)*len(inp.v[0])/8+1)
	bits := 0
	bbuf := byte(0)
	bbufpos := byte(0)
	ptr := 0
	for row, v := range inp.v {
		for col := range v {
			if inp.marks[pos{row: row + offsetrow, col: col + offsetcol}] {
				bbuf = bbuf | (1 << bbufpos)
				bits++
			}
			bbufpos++
			if bbufpos >= 8 {
				bbufpos = 0
				res[ptr] = bbuf
				bbuf = 0
				ptr++
			}
		}
	}
	res[ptr] = bbuf
	return &bitmap{
		buff: res,
		bits: bits,
	}
}

func (inp *Input) AllBitMaps() map[pos]*bitmap {
	res := map[pos]*bitmap{}
	rows := len(inp.m)
	cols := len(inp.m[0])
	for offsetrow := inp.min.row - modNeg(inp.min.row, rows); offsetrow <= inp.max.row; offsetrow += rows {
		for offsetcol := inp.min.col - modNeg(inp.min.col, cols); offsetcol <= inp.max.col; offsetcol += cols {
			res[pos{row: offsetrow, col: offsetcol}] = inp.BitMap(offsetrow, offsetcol)
		}
	}
	return res
}

func (inp *Input) AllBitMapsWithNeighbors() map[pos]*[]byte {
	res := map[pos]*[]byte{}
	rows := len(inp.m)
	cols := len(inp.m[0])
	for offsetrow := inp.min.row - modNeg(inp.min.row, rows); offsetrow < inp.max.row; offsetrow += rows {
		for offsetcol := inp.min.col - modNeg(inp.min.col, cols); offsetcol < inp.max.col; offsetcol += cols {
			res[pos{row: offsetrow, col: offsetcol}] = inp.BitMapWithNeighbors(offsetrow, offsetcol)
		}
	}
	return res
}

func (inp *Input) BitMapWithNeighbors(offsetrow, offsetcol int) *[]byte {
	rows := len(inp.v)
	cols := len(inp.v[0])
	res := make([]byte, (rows*cols+2*rows+2*cols)/8+1)
	bbuf := byte(0)
	bbufpos := byte(0)
	ptr := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if inp.marks[pos{row: row + offsetrow, col: col + offsetcol}] {
				bbuf = bbuf | (1 << bbufpos)
			}
			bbufpos++
			if bbufpos >= 8 {
				bbufpos = 0
				res[ptr] = bbuf
				bbuf = 0
				ptr++
			}
		}
	}
	for row := 0; row < rows; row++ {
		if inp.marks[pos{row: row + offsetrow, col: offsetcol - 1}] {
			bbuf = bbuf | (1 << bbufpos)
		}
		bbufpos++
		if bbufpos >= 8 {
			bbufpos = 0
			res[ptr] = bbuf
			bbuf = 0
			ptr++
		}
		if inp.marks[pos{row: row + offsetrow, col: offsetcol + cols}] {
			bbuf = bbuf | (1 << bbufpos)
		}
		bbufpos++
		if bbufpos >= 8 {
			bbufpos = 0
			res[ptr] = bbuf
			bbuf = 0
			ptr++
		}
	}
	for col := 0; col < cols; col++ {
		if inp.marks[pos{row: offsetrow - 1, col: col + offsetcol}] {
			bbuf = bbuf | (1 << bbufpos)
		}
		bbufpos++
		if bbufpos >= 8 {
			bbufpos = 0
			res[ptr] = bbuf
			bbuf = 0
			ptr++
		}
		if inp.marks[pos{row: offsetrow + rows, col: col + offsetcol}] {
			bbuf = bbuf | (1 << bbufpos)
		}
		bbufpos++
		if bbufpos >= 8 {
			bbufpos = 0
			res[ptr] = bbuf
			bbuf = 0
			ptr++
		}
	}
	return &res
}

type dir struct {
	row, col int
}

var ltrs = map[byte]dir{
	'N': {
		row: -1,
		col: 0,
	},
	'S': {
		row: 1,
		col: 0,
	},
	'E': {
		row: 0,
		col: 1,
	},
	'W': {
		row: 0,
		col: -1,
	},
}

var alldirs = []dir{
	{
		row: -1,
		col: 0,
	},
	{
		row: 1,
		col: 0,
	},
	{
		row: 0,
		col: -1,
	},
	{
		row: 0,
		col: 1,
	},
}

func (a pos) Add(b dir) pos {
	return pos{row: a.row + b.row, col: a.col + b.col}
}

func (a dir) Opposed(b dir) bool {
	return a.row == -b.row && a.col == -b.col
}

func (a dir) Reverse() dir {
	return dir{row: -a.row, col: -a.col}
}

func (inp Input) Print() {
	for row, v := range inp.v {
		for col, v2 := range v {
			if v2 {
				if inp.marks[pos{row: row, col: col}] {
					fmt.Print("O")
				} else {
					fmt.Print(".")
				}
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (s pos) Neighbors(inp *Input) []pos {
	res := []pos{}
	for _, d := range ltrs {
		s1 := s.Add(d)
		if s1.col < 0 || s1.row < 0 || s1.col >= len(inp.m[0]) || s1.row >= len(inp.m) {
			continue
		}
		if inp.v[s1.row][s1.col] {
			res = append(res, s1)
		}
	}
	return res
}

func modNeg(v, m int) int {
	return (v%m + m) % m
}
func (s pos) Neighbours2(inp Input) []pos {
	res := []pos{}
	for _, d := range ltrs {
		s1 := s.Add(d)
		row := modNeg(s1.row, len(inp.m))
		col := modNeg(s1.col, len(inp.m[0]))
		if inp.v[row][col] {
			res = append(res, s1)
		}
	}
	return res
}

func (inp Input) Count(steps int, steps2 int) (int, int) {
	s1 := 0
	s2 := 0
	startingmarks := inp.marks
	for i := 0; i < steps; i++ {
		m := map[pos]bool{}
		for p := range inp.marks {
			for _, p2 := range p.Neighbors(&inp) {
				m[p2] = true
			}
		}
		inp.marks = m
		//		inp.Print()
	}
	s1 = len(inp.marks)
	inp.marks = startingmarks
	bmhash := map[string]map[int]int{}
	spreds := 0
	for i := 0; i < steps2; i++ {
		m := map[pos]bool{}
		//bmsstart := inp.AllBitMapsWithNeighbors()
		for p := range inp.marks {
			for _, p2 := range p.Neighbours2(inp) {
				m[p2] = true
				inp.min.col = min(inp.min.col, p2.col)
				inp.min.row = min(inp.min.row, p2.row)
				inp.max.col = max(inp.max.col, p2.col)
				inp.max.row = max(inp.max.row, p2.row)
			}
		}
		inp.marks = m
		cs := 0
		bms := inp.AllBitMaps()
		for _, v := range bms {
			sv := v.String()
			if _, ok := bmhash[sv]; !ok {
				bmhash[sv] = map[int]int{i: 1}
			} else {
				bmhash[sv][i]++
			}
			cs += v.bits
			/*			if v.bits != bitMapFromString(sv).bits {
						fmt.Printf("?????")
					}*/
		}
		if cs != len(inp.marks) {
			fmt.Printf("!!!!!")
		}
		//		inp.Print()
		//fmt.Println(i, " ", len(m), " ", len(bmhash))
		s2exp := inp.Predict(bmhash, i+1)
		if len(inp.marks) == s2exp {
			fmt.Printf("!!!!! %d\n", i)
			spreds += 1
		}
		if s2exp != 0 {
			fmt.Println("--- ", len(inp.marks), " ", s2exp, " +", spreds)
		}
		if spreds > 2*len(inp.m) {
			return s1, inp.Predict(bmhash, steps2)
		}
	}
	s2exp := inp.Predict(bmhash, steps)
	s2 = len(inp.marks)
	fmt.Printf("s2: %d s2exp: %d", s2, s2exp)
	return s1, s2
}

func (inp Input) Predict(bmhash map[string]map[int]int, steps int) int {
	s2exp := 0
	for bms, v := range bmhash {

		if len(v) > 1 {
			expInterval := len(inp.m)
			series, err := ExamineMap(v, expInterval)
			if err == nil {
				for _, sd := range series {
					exemplars := sd.Predict(steps-1, expInterval)
					if exemplars > 0 {
						s2exp += exemplars * bitMapFromString(bms).bits
					}
				}
			}
		}
	}
	return s2exp
}

type seriesDesc struct {
	basei     int
	basev     int
	increment int
	quadratic int
}

func (s seriesDesc) Predict(time, interval int) int {
	if time%interval == s.basei%interval {
		steps := (time - s.basei) / interval
		if s.increment > 1 {
			fmt.Println("inc:", s.increment)
		}
		return s.basev + steps*s.increment + (steps*(steps-1)*s.quadratic)/2
	}
	return 0
}

func ExamineMap(v map[int]int, expinterval int) ([]seriesDesc, error) {
	keys := make([]int, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	if len(keys) < 4 {
		return []seriesDesc{}, fmt.Errorf("Not enough keys")
	}
	res := []seriesDesc{}
	refkey := keys[len(keys)-1] - 2*expinterval
	refkey1 := keys[len(keys)-1] - expinterval
	refkey2 := refkey - expinterval
	if value, ok := v[refkey]; ok {
		if value1, ok := v[refkey1]; ok {
			quadratic := 0
			inc1 := value1 - value
			if value2, ok := v[refkey2]; ok {
				inc2 := value - value2
				quadratic = inc1 - inc2
			}
			res = append(res, seriesDesc{
				basei:     refkey,
				basev:     value,
				increment: inc1,
				quadratic: quadratic,
			})
		} else {
			return []seriesDesc{}, fmt.Errorf("Varying interval1")
		}
	} else {
		return []seriesDesc{}, fmt.Errorf("Varying interval")
	}
	for idx := len(keys) - 2; idx > 0 && keys[idx]%expinterval != keys[len(keys)-1]%expinterval; idx-- {
		refkey := keys[idx] - 2*expinterval
		refkey1 := keys[idx] - expinterval
		refkey2 := refkey - expinterval
		if value, ok := v[refkey]; ok {
			if value1, ok := v[refkey1]; ok {
				quadratic := 0
				inc1 := value1 - value
				if value2, ok := v[refkey2]; ok {
					inc2 := value - value2
					quadratic = inc1 - inc2
				}
				res = append(res, seriesDesc{
					basei:     refkey,
					basev:     value,
					increment: inc1,
					quadratic: quadratic,
				})
			}
		}
	}
	return res, nil
}

func Readlines(file io.Reader) (int, int) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	return inp.Count(64, 26501365)
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
