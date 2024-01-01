package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type Line struct {
	winningbyte [100]bool
	hand        []int
	number      int
}
type Copies map[int]int

func (c Copies) win(current, score int) {
	for p := current + 1; p <= current+score; p++ {
		if _, ok := c[p]; !ok {
			c[p] = c[current]
		} else {
			c[p] += c[current]
		}
	}
}

func (c Copies) take(current int) {
	if _, ok := c[current]; !ok {
		c[current] = 1
	} else {
		c[current] += 1
	}
}

func (c Copies) count() int {
	res := 0
	for _, v := range c {
		res += v
	}
	return res
}

var lineRe = regexp.MustCompile(`Card *([0-9]+): (.*) \| (.*)`)

var intmap = map[string]int{
	" 1": 1,
	" 2": 2,
	" 3": 3,
	" 4": 4,
	" 5": 5,
	" 6": 6,
	" 7": 7,
	" 8": 8,
	" 9": 9,
	"10": 10,
	"11": 11,
	"12": 12,
	"13": 13,
	"14": 14,
	"15": 15,
	"16": 16,
	"17": 17,
	"18": 18,
	"19": 19,
	"20": 20,
	"21": 21,
	"22": 22,
	"23": 23,
	"24": 24,
	"25": 25,
	"26": 26,
	"27": 27,
	"28": 28,
	"29": 29,
	"30": 30,
	"31": 31,
	"32": 32,
	"33": 33,
	"34": 34,
	"35": 35,
	"36": 36,
	"37": 37,
	"38": 38,
	"39": 39,
	"40": 40,
	"41": 41,
	"42": 42,
	"43": 43,
	"44": 44,
	"45": 45,
	"46": 46,
	"47": 47,
	"48": 48,
	"49": 49,
	"50": 50,
	"51": 51,
	"52": 52,
	"53": 53,
	"54": 54,
	"55": 55,
	"56": 56,
	"57": 57,
	"58": 58,
	"59": 59,
	"60": 60,
	"61": 61,
	"62": 62,
	"63": 63,
	"64": 64,
	"65": 65,
	"66": 66,
	"67": 67,
	"68": 68,
	"69": 69,
	"70": 70,
	"71": 71,
	"72": 72,
	"73": 73,
	"74": 74,
	"75": 75,
	"76": 76,
	"77": 77,
	"78": 78,
	"79": 79,
	"80": 80,
	"81": 81,
	"82": 82,
	"83": 83,
	"84": 84,
	"85": 85,
	"86": 86,
	"87": 87,
	"88": 88,
	"89": 89,
	"90": 90,
	"91": 91,
	"92": 92,
	"93": 93,
	"94": 94,
	"95": 95,
	"96": 96,
	"97": 97,
	"98": 98,
	"99": 99,
}

func parseNumberList(s string) []int {
	res := make([]int, (len(s)+1)/3)
	p1 := 0
	for p := 0; p < len(s); p += 3 {
		res[p1] = intmap[s[p:p+2]]
		p1 += 1
	}
	return res
}
func ParseLine(s string) (Line, error) {
	l := Line{
		winningbyte: [100]bool{},
		number:      0,
	}
	r := lineRe.FindStringSubmatchIndex(s)
	if r == nil {
		return l, fmt.Errorf("parse error for %s", s)
	}
	l.number, _ = strconv.Atoi(s[r[2]:r[3]])

	for _, v := range parseNumberList(s[r[4]:r[5]]) {
		l.winningbyte[v] = true
	}
	/*ws := strings.Split(s[r[4]:r[5]], " ")
	for _, wss := range ws {
		if wsi, err := strconv.Atoi(wss); err == nil {
			l.winningbyte[wsi] = true
		}
	}*/
	l.hand = parseNumberList(s[r[6]:r[7]])
	/*hs := strings.Split(s[r[6]:r[7]], " ")
	for _, hss := range hs {
		if hsi, err := strconv.Atoi(hss); err == nil {
			l.hand = append(l.hand, hsi)
		}
	}*/
	return l, nil
}

func (l *Line) Score() int {
	s := 0
	for _, v := range l.hand {
		if l.winningbyte[v] {
			if s == 0 {
				s = 1
				continue
			}
			s = s * 2
		}
	}
	return s
}

func (l *Line) Matches() int {
	s := 0
	for _, v := range l.hand {
		if l.winningbyte[v] {
			s += 1
		}
	}
	return s
}

func Readlines(file io.Reader) (int, int) {
	res := 0
	copies := Copies{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s, e := ParseLine(scanner.Text())
		if e == nil {
			res += s.Score()
			copies.take(s.number)
			copies.win(s.number, s.Matches())
		}
	}
	return res, copies.count()
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		runtime.SetCPUProfileRate(100000)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Error opening input file: ", err)
	}
	defer file.Close()
	n := time.Now()
	f, f1 := Readlines(file)
	fmt.Println(strconv.Itoa(f), strconv.Itoa(f1))
	fmt.Println(time.Since(n))
}
