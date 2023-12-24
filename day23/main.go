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
	m     []string
	v     map[pos]*vertex
	e     []*edge
	start pos
	end   pos
}

type pos dir
type vertexType byte

type vertex struct {
	p pos
	e map[dir]*edge
	t vertexType
	d dir
}

const (
	vtDir     = vertexType(1)
	vtBranch  = vertexType(2)
	vtStart   = vertexType(3)
	vtEnd     = vertexType(4)
	vtDeadEnd = vertexType(5)
)

type edge struct {
	d1     dir
	d2     dir
	v1     *vertex
	v2     *vertex
	length int
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m: []string{},
		v: map[pos]*vertex{},
		e: []*edge{},
	}
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		row++
	}
	pstart := pos{row: 0, col: 1}
	res.start = pstart
	res.v[pstart] = &vertex{
		p: pstart,
		e: map[dir]*edge{},
		t: vtStart,
		d: S,
	}
	pend := pos{
		row: row - 1,
		col: len(res.m[row-1]) - 2,
	}
	res.end = pend
	res.v[pend] = &vertex{
		p: pend,
		e: map[dir]*edge{},
		t: vtEnd,
		d: N,
	}
	res.ScanGraph(res.v[pstart], S)
	return res
}

type dir struct {
	row, col int
}

var (
	N       = dir{row: -1, col: 0}
	S       = dir{row: 1, col: 0}
	E       = dir{row: 0, col: 1}
	W       = dir{row: 0, col: -1}
	DIRS    = []dir{N, S, E, W}
	CHTODIR = map[byte]dir{
		'>': E,
		'<': W,
		'v': S,
		'^': N,
	}
)

func (a pos) Add(b dir) pos {
	return pos{row: a.row + b.row, col: a.col + b.col}
}

func (a dir) Opposed(b dir) bool {
	return a.row == -b.row && a.col == -b.col
}

func (a dir) Reverse() dir {
	return dir{row: -a.row, col: -a.col}
}

func (inp *Input) IsWall(p pos) bool {
	if p.row < 0 || p.row >= len(inp.m) || p.col < 0 || p.col > len(inp.m[p.row]) {
		return false
	}
	return inp.m[p.row][p.col] == '#'
}

func (inp *Input) IsVertex(p pos) (bool, vertexType) {
	switch inp.m[p.row][p.col] {
	case '#':
		return false, 0
	case '>', '<', 'v', '^':
		return true, vtDir
	default:
		nbr := 0
		if _, ok := inp.v[p]; ok {
			return true, inp.v[p].t
		}
		for _, d := range DIRS {
			p1 := p.Add(d)
			if inp.m[p1.row][p1.col] != '#' {
				nbr++
			}
		}
		switch nbr {
		case 1:
			return true, vtDeadEnd
		case 2:
			return false, 0
		default:
			return true, vtBranch
		}
	}
}

func (inp *Input) ScanGraph(v *vertex, d dir) {
	if inp.IsWall(v.p.Add(d)) {
		return
	}
	currentpos := v.p
	currentdir := d
	currentEdge := edge{
		d1:     d,
		d2:     d.Reverse(),
		v1:     v,
		v2:     v,
		length: 0,
	}
	v.e[currentdir] = &currentEdge
	for {
		nextpos := currentpos.Add(currentdir)
		currentEdge.length++
		if ok, vt := inp.IsVertex(nextpos); ok {
			if v2, ok := inp.v[nextpos]; ok {
				currentEdge.v2 = v2
				currentEdge.d2 = currentdir.Reverse()
				v2.e[currentdir.Reverse()] = &currentEdge
				inp.e = append(inp.e, &currentEdge)
				return
			} else {
				d2 := dir{0, 0}
				if d2v, ok := CHTODIR[inp.m[nextpos.row][nextpos.col]]; ok {
					d2 = d2v
				}
				vnew := vertex{
					p: nextpos,
					e: map[dir]*edge{currentdir.Reverse(): &currentEdge},
					t: vt,
					d: d2,
				}
				inp.v[nextpos] = &vnew
				currentEdge.v2 = &vnew
				inp.e = append(inp.e, &currentEdge)
				for _, d3 := range DIRS {
					if d3.Reverse() != currentdir {
						inp.ScanGraph(&vnew, d3)
					}
				}
				return
			}
		} else {
			for _, d2 := range DIRS {
				if d2.Reverse() != currentdir {
					if !inp.IsWall(nextpos.Add(d2)) {
						currentpos = nextpos
						currentdir = d2
						break
					}
				}
			}
		}
	}
}
func (inp Input) Has(point pos) (byte, bool) {
	if point.row < 0 || point.col < 0 {
		return 0, false
	}
	if point.row >= len(inp.m) {
		return 0, false
	}
	if point.col >= len(inp.m[point.row]) {
		return 0, false
	}
	return inp.m[point.row][point.col], true
}

func (e *edge) OtherVertexDir(v *vertex) (*vertex, dir) {
	if v == e.v1 {
		return e.v2, e.d2
	} else {
		return e.v1, e.d1
	}
}

func (inp Input) LongestPath(visited map[pos]*vertex, lastvertex *vertex, currentlength int) int {
	end := inp.v[inp.end]
	maxdistance := 0
	for d, e := range lastvertex.e {
		if lastvertex.t == vtDir && d != lastvertex.d {
			continue
		}
		if e.v1 == end {
			maxdistance = max(e.length, maxdistance)
			continue
		}
		if e.v2 == end {
			maxdistance = max(e.length, maxdistance)
			continue
		}
		ov, _ := e.OtherVertexDir(lastvertex)
		if _, ok := visited[ov.p]; ok {
			continue
		}
		v1 := map[pos]*vertex{}
		v1[ov.p] = ov
		for p, v := range visited {
			v1[p] = v
		}
		lp := inp.LongestPath(v1, ov, currentlength+e.length) - currentlength
		maxdistance = max(lp, maxdistance)
	}
	if maxdistance > 0 {
		return maxdistance + currentlength
	}
	return 0
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	s1 = inp.LongestPath(map[pos]*vertex{inp.start: inp.v[inp.start]}, inp.v[inp.start], 0)
	for i := range inp.m {
		inp.m[i] = strings.Replace(inp.m[i], ">", ".", -1)
		inp.m[i] = strings.Replace(inp.m[i], "<", ".", -1)
		inp.m[i] = strings.Replace(inp.m[i], "v", ".", -1)
		inp.m[i] = strings.Replace(inp.m[i], "^", ".", -1)
	}
	inp.e = []*edge{}
	inp.v = map[pos]*vertex{
		inp.start: {
			p: inp.start,
			e: map[dir]*edge{},
			t: vtStart,
			d: S,
		},
		inp.end: {
			p: inp.end,
			e: map[dir]*edge{},
			t: vtEnd,
			d: N,
		},
	}
	inp.ScanGraph(inp.v[inp.start], S)
	s2 = inp.LongestPath(map[pos]*vertex{inp.start: inp.v[inp.start]}, inp.v[inp.start], 0)

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
