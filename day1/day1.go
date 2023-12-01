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

func getCalibrationValue(s string) int {
	digits := regexp.MustCompile("[0-9]")
	v1 := ""
	v2 := ""
	for _, v := range strings.Split(s, "") {
		if digits.Match([]byte(v)) {
			if v1 == "" {
				v1 = v
			}
			v2 = v
		}
	}
	res, _ := strconv.Atoi(v1 + v2)
	return res
}

var names = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

func getCalibrationValue2(s string, rs []*regexp.Regexp) int {
	v1 := -1
	v2 := 0
	for i, _ := range strings.Split(s, "") {
		vl := s[i:min(i+5, len(s))]
		for i2, r := range rs {
			if r.Match([]byte(vl)) {
				if v1 == -1 {
					v1 = i2
				}
				v2 = i2
			}
		}
	}
	if v1 == -1 {
		return 0
	}
	return v2 + 10*v1
}
func Readlines(file io.Reader) int {
	regexps := []*regexp.Regexp{regexp.MustCompile("^0")}
	for i, v := range names {
		regexps = append(regexps, regexp.MustCompile(
			"^"+strconv.Itoa(i+1)+"|^"+v,
		))
	}

	scanner := bufio.NewScanner(file)
	res := 0
	res2 := 0
	for scanner.Scan() {
		s := scanner.Text()
		res += getCalibrationValue(s)

		cv2 := getCalibrationValue2(s, regexps)
		fmt.Println(s, cv2)
		res2 += cv2
	}
	fmt.Println(res)
	fmt.Println(res2)
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
