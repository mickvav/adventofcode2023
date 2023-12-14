package main

import (
	"bufio"
	"fmt"
	"hash/maphash"
	"io"
	"log"
	"os"
	"strconv"
)

type Input struct {
	m       []string
	b       [][]byte
	h       maphash.Seed
	history map[uint64]int
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:       []string{},
		b:       [][]byte{},
		h:       maphash.MakeSeed(),
		history: map[uint64]int{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		res.b = append(res.b, []byte(line))
	}
	return res
}

func (inp *Input) tiltNorth() {
	moved := true
	for moved {
		moved = false
		for i, v := range inp.b {
			if i > 0 {
				for i2, v2 := range v {
					if v2 == 'O' && inp.b[i-1][i2] == '.' {
						moved = true
						inp.b[i-1][i2] = 'O'
						inp.b[i][i2] = '.'
					}
				}
			}
		}
	}
}

func (inp *Input) tiltSouth() {
	moved := true
	for moved {
		moved = false
		for i, v := range inp.b {
			if i < len(inp.b)-1 {
				for i2, v2 := range v {
					if v2 == 'O' && inp.b[i+1][i2] == '.' {
						moved = true
						inp.b[i+1][i2] = 'O'
						inp.b[i][i2] = '.'
					}
				}
			}
		}
	}
}

func (inp *Input) tiltEast() {
	moved := true
	for moved {
		moved = false
		for i, v := range inp.b {
			for i2, _ := range v {
				if i2 < len(v)-1 {
					if inp.b[i][i2] == 'O' && inp.b[i][i2+1] == '.' {
						moved = true
						inp.b[i][i2+1] = 'O'
						inp.b[i][i2] = '.'
					}
				}
			}
		}
	}
}

func (inp *Input) tiltWest() {
	moved := true
	for moved {
		moved = false
		for i, v := range inp.b {
			for i2, _ := range v {
				if i2 > 0 {
					if inp.b[i][i2] == 'O' && inp.b[i][i2-1] == '.' {
						moved = true
						inp.b[i][i2-1] = 'O'
						inp.b[i][i2] = '.'
					}
				}
			}
		}
	}
}

func (inp *Input) Hash() uint64 {
	var h maphash.Hash
	h.SetSeed(inp.h)
	h.Reset()
	for _, v := range inp.b {
		h.Write(v)
	}
	return h.Sum64()
}

func (inp *Input) Print() {
	for _, v := range inp.b {
		fmt.Printf(string(v) + "\n")
	}
}
func (inp *Input) Count() (int, int) {
	res := 0
	rows := len(inp.b)
	for _, v := range inp.b {
		for _, v2 := range v {
			if v2 == 'O' {
				res += rows
			}
		}
		rows = rows - 1
	}
	return res, 0
}

func Readlines(file io.Reader) (int, int) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	inp.tiltNorth()
	n1, n2 := inp.Count()
	for j := 1; j < 10000; j++ {
		inp.tiltNorth()
		inp.tiltWest()
		inp.tiltSouth()
		inp.tiltEast()
		n2, _ = inp.Count()
		fmt.Println(j, n2)
		h := inp.Hash()
		if s, ok := inp.history[h]; ok {
			if 1000000000%(j-s) == s {
				return n1, n2
			}
		} else {
			inp.history[h] = j
		}
	}
	return n1, n2
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
