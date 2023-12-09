package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)


type Input struct {
	instructions string
	L map[string]string
	R map[string]string
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		instructions: "",
		L:            map[string]string{},
		R:            map[string]string{},
	}
	scanner.Scan()
	res.instructions = scanner.Text()
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			k := line[0:3]
			res.L[k] = line[7:10]
			res.R[k] = line[12:15]
		}
	}
	return res
}

func (inp Input) Count() int {
	s := 0
	p := 0
	k := "AAA"
	for k != "ZZZ" {
		cmd := inp.instructions[p]
		if cmd == 'L' {
			k=inp.L[k]
		}
		if cmd == 'R' {
			k=inp.R[k]
		}
		s=s+1
		p=(p+1) % len(inp.instructions)
	}
	return s
}

func (inp Input) Starts() []string {
	res := []string{}
	for k, _ := range inp.L {
		if k[2] == 'A' {
			res = append(res, k)
		}
	}
	return res
}

type data struct {
	i1z map[int]int
	i2z map[int]int
}
func (d data)Finish(keys *[]string, iter int) bool {
	res := true
	if keys != nil {
		for g, v := range (*keys) {
			if v[2] != 'Z' {
				res = false
			} else {
				if _,ok := d.i1z[g]; !ok {
					d.i1z[g] = iter
					continue
				}

				if _,ok := d.i2z[g]; !ok {
					d.i2z[g] = iter
					continue
				}
			}
		}
	}
	return len(d.i2z) == len(*keys) || res
}
func (inp Input) Count2() int64 {
	s := 0
	p := 0
	k := inp.Starts()
	d:=data{
		i1z: map[int]int{},
		i2z: map[int]int{},
	}
	for !d.Finish(&k, s) {
		for i, v := range k {
			if v[2] == 'Z' {
				fmt.Println(i, s, v)
			}
		}
		cmd := inp.instructions[p]
		if cmd == 'L' {
			for i, v := range k {
				k[i] = inp.L[v]
			}
		}
		if cmd == 'R' {
			for i, v := range k {
				k[i] = inp.R[v]
			}
		}
		s=s+1
		p=(p+1) % len(inp.instructions)
	}
	c := []int64{}
	for _, v := range d.i1z {
		c = append(c, int64(v))
	}
	return LCM(c[0], c[1], c[2:]...)
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int64, integers ...int64) int64 {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func Readlines(file io.Reader) (int, int64) {
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
	fmt.Println(strconv.Itoa(f), f1)
}
