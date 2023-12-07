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
	times     []int
	distances []int
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		times:     []int{},
		distances: []int{},
	}
	scanner.Scan()
	line1 := scanner.Text()
	times := strings.Split(line1, " ")
	for _, s := range times[1:] {
		if s != "" {
			sv, _ := strconv.Atoi(s)
			res.times = append(res.times, sv)
		}
	}
	scanner.Scan()
	line1 = scanner.Text()
	distances := strings.Split(line1, " ")
	for _, s := range distances[1:] {
		if s != "" {
			sv, _ := strconv.Atoi(s)
			res.distances = append(res.distances, sv)
		}
	}
	return res
}

func (inp Input) Join() Input {
	res := Input{
		times:     []int{},
		distances: []int{},
	}
	ds :=""
	for _, v := range inp.distances {
		ds += strconv.Itoa(v)
	}
	di, _ := strconv.Atoi(ds)
	res.distances = append(res.distances, di)
	ts := ""
	for _, v := range inp.times {
		ts += strconv.Itoa(v)
	}
	ti, _ := strconv.Atoi(ts)
	res.times = append(res.times, ti)
	return res
}

func (inp Input) CountVariants() int {
	mult := 1
	for i := 0; i < len(inp.distances); i++ {
		res := 0
		for tw := 0; tw < inp.times[i]; tw++ {
			v := tw
			d := v * (inp.times[i] - tw)
			if d > inp.distances[i] {
				res += 1
			}
		}
		mult = mult * res
	}
	return mult
}

func Readlines(file io.Reader) (int, int) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	return inp.CountVariants(), inp.Join().CountVariants()
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
