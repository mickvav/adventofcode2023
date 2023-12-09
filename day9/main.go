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
	histories [][]int
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		histories: [][]int{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			values := []int{}
			stringvalues := strings.Split(line, " ")
			for _, v := range stringvalues {
				value, _ := strconv.Atoi(v)
				values = append(values, value)
			}
			res.histories = append(res.histories, values)
		}
	}
	return res
}

func Predict(inp []int) (int, int) {
	diffs := [][]int{}
	finish := false
	d := make([]int, len(inp))
	copy(d,inp)
	for !finish {
		d, finish = Diffs(d)
		diffs = append(diffs, d)
	}
	prevdiff := 0
	prevdiffStart :=0 
	for l := len(diffs) -1 ; l>=0 ; l-- {
		mld := len(diffs[l])-1
		diffs[l] = append(diffs[l], diffs[l][mld] + prevdiff )
		prevdiff = diffs[l][mld+1]
		prevdiffStart = diffs[l][0] - prevdiffStart
	}
	return inp[len(inp)-1] + prevdiff, inp[0] - prevdiffStart
}


func Diffs(inp []int) (res []int, finish bool)  {
	finish = true
	for i, v := range inp {
		if i > 0 {
			d := v - inp[i-1]
			res = append(res, d)
			if d != 0 {
				finish = false
			}
		}
	}
	return
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	for _, v := range inp.histories {
		p1, p2 := Predict(v)
		s1 += p1
		s2 += p2
	}
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
