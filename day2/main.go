package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type game struct {
	maxRed   int
	maxGreen int
	maxBlue  int
	number   int
}

// Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
// Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
// Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
// Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
// Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
var Line = regexp.MustCompile("Game ([0-9]+):(.*)")
var Lb = regexp.MustCompile("([0-9]+) (red|green|blue)")

func ParseLine(s string) (game, error) {
	res := game{
		maxRed:   0,
		maxGreen: 0,
		maxBlue:  0,
		number:   0,
	}
	m := Line.FindStringSubmatch(s)
	if m == nil {
		return res, fmt.Errorf("Parse error")
	}
	var err error
	res.number, err = strconv.Atoi(m[1])
	if err != nil {
		return res, err
	}
	rounds := strings.Split(m[2], "; ")
	for _, v := range rounds {
		bunches := strings.Split(v, ", ")
		for _, v2 := range bunches {
			m1 := Lb.FindStringSubmatch(v2)
			if m1 == nil {
				return res, fmt.Errorf("Cant parse %s", v2)
			}
			n, _ := strconv.Atoi(m1[1])
			switch m1[2] {
			case "red":
				res.maxRed = max(res.maxRed, n)
			case "green":
				res.maxGreen = max(res.maxGreen, n)
			case "blue":
				res.maxBlue = max(res.maxBlue, n)
			}
		}
	}
	return res, nil
}

func Readlines(file io.Reader) int {
	scanner := bufio.NewScanner(file)
	res := 0
	res1 := 0
	for scanner.Scan() {
		s := scanner.Text()
		g, err := ParseLine(s)
		if err != nil {
			fmt.Print(err)
		}
		if g.maxRed <= 12 && g.maxBlue <= 14 && g.maxGreen <= 13 {
			// possible
			res += g.number
		}
		res1 += g.maxRed * g.maxBlue * g.maxGreen
	}
	fmt.Println(res)
	fmt.Println(res1)
	return res
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Ups")
	}
	defer file.Close()
	fmt.Println(strconv.Itoa(Readlines(file)))
}
