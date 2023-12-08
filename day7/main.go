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

type hand struct {
	repr string
	v    [5]int
	hist map[int]int
	rank int
	rank2 int
	v2   [5]int
	bid  int
}

func parsehand(s string) hand {
	h := hand{
		repr: s,
		v:    [5]int{},
		hist: map[int]int{},
		rank: 0,
		bid:  0,
	}
	p := strings.Split(s, " ")
	h.bid, _ = strconv.Atoi(p[1])
	for i, v := range strings.Split(p[0], "") {
		switch v {
		case "2":
			h.v[i] = 2
		case "3":
			h.v[i] = 3
		case "4":
			h.v[i] = 4
		case "5":
			h.v[i] = 5
		case "6":
			h.v[i] = 6
		case "7":
			h.v[i] = 7
		case "8":
			h.v[i] = 8
		case "9":
			h.v[i] = 9
		case "T":
			h.v[i] = 10
		case "J":
			h.v[i] = 11
		case "Q":
			h.v[i] = 12
		case "K":
			h.v[i] = 13
		case "A":
			h.v[i] = 14
		}
		h.v2[i] = h.v[i]
		if h.v[i]== 11 {
			h.v2[i] = 1 // J
		}
		h.hist[h.v[i]]++
	}
	if len(h.hist) == 1 {
		h.rank = 7
		h.rank2 = 7
	}
	if len(h.hist) == 2 {
		h.rank = 5
		for _, v := range h.hist {
			if v == 1 {
				h.rank = 6
				break
			}
		}
		h.rank2 = h.rank
		if _, ok := h.hist[11]; ok {
			h.rank2 = 7
		}
	}
	if len(h.hist) == 3 {
		h.rank = 3
		for _, v := range h.hist {
			if v == 3 {
				h.rank = 4
			}
		}
		h.rank2 = h.rank
		if v, ok := h.hist[11]; ok {
			if v == 3 { // 3J -> 4 of kind
				h.rank2 = 6
			}
			if v == 2 { // 2J -> 4 of kind
				h.rank2 = 6
			}
			if v == 1 {
				if h.rank == 3 {// 1J + two pair - full house 
					h.rank2 = 5
				}
				if h.rank == 4 {// 1J + 3 of kind -> 4 of kind
					h.rank2 = 6
				}
			}
		}
	}
	if len(h.hist) == 4 {
		h.rank = 2
		h.rank2 = 2
		if _, ok := h.hist[11]; ok {
			h.rank2 = 4 
		}
	}
	if len(h.hist) == 5 {
		h.rank = 1
		h.rank2 = 1
		if _, ok := h.hist[11]; ok {
			h.rank2 = 2
		}
	}
	return h
}

type ByOrder []hand
type ByOrder2 []hand

func (a ByOrder) Len() int      { return len(a) }
func (a ByOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByOrder) Less(i, j int) bool {
	if a[i].rank != a[j].rank {
		return a[i].rank < a[j].rank
	}
	for k := 0; k < 5; k++ {
		if a[i].v[k] != a[j].v[k] {
			return a[i].v[k] < a[j].v[k]
		}
	}
	return false
}
func (a ByOrder2) Len() int      { return len(a) }
func (a ByOrder2) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByOrder2) Less(i, j int) bool {
	if a[i].rank2 != a[j].rank2 {
		return a[i].rank2 < a[j].rank2
	}
	for k := 0; k < 5; k++ {
		if a[i].v2[k] != a[j].v2[k] {
			return a[i].v2[k] < a[j].v2[k]
		}
	}
	return false
}


type Input struct {
	hands []hand
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		hands: []hand{},
	}
	for scanner.Scan() {
		line1 := scanner.Text()
		res.hands = append(res.hands, parsehand(line1))
	}
	return res
}

func (inp Input) Count() int {
	sort.Sort(ByOrder(inp.hands))
	res := 0
	for i, h := range inp.hands {
		res += (i + 1) * h.bid
	}
	return res
}
func (inp Input) Count2() int {
	sort.Sort(ByOrder2(inp.hands))
	res := 0
	for i, h := range inp.hands {
		res += (i + 1) * h.bid
	}
	return res
}



func Readlines(file io.Reader) (int, int) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	return inp.Count(), inp.Count2()
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
