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

type Interval struct {
	source int
	destination int
	length int
}
type Ranger struct {
	intervals []Interval
}
type BySource []Interval 


func (a BySource) Len() int           { return len(a) }
func (a BySource) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySource) Less(i, j int) bool { return a[i].source < a[j].source }

func ReadRanger(scaner *bufio.Scanner) Ranger {
	r := Ranger{
		intervals: []Interval{},
	}
	l:="==="
	for scaner.Scan() && l!= "" {
		l = scaner.Text()
		tokens := strings.Split(l, " ")
		if len(tokens) == 3 {
			dst, _ := strconv.Atoi(tokens[0])
			src, _ := strconv.Atoi(tokens[1])
			length, _ := strconv.Atoi(tokens[2])
			r.intervals = append(r.intervals, Interval{
				source: src, destination: dst, length: length,
			})
		}
	}
	sort.Sort(BySource(r.intervals))
	return r
}

func (r Ranger) FindDest(src int) int {
	for _, i2 := range r.intervals {
		if src < i2.source {
			return src
		}
		if src <= i2.source + i2.length {
			return i2.destination + src - i2.source
		}
	}
	return src
}

type Input struct {
	Seeds []int
	SeedSoilMap Ranger
	SoilFertMap Ranger
	FertWaterMap Ranger
	WaterLightMap Ranger
	LightTempMap Ranger
	TemphumMap Ranger
	HumLocMap Ranger
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		Seeds:         []int{},
	}
	scanner.Scan()
	line1 := scanner.Text()
	seeds := strings.Split(line1, " ")
	for _, s := range(seeds[1:]) {
		sv, _ := strconv.Atoi(s)
		res.Seeds = append(res.Seeds, sv)
	}
	scanner.Scan()
	res.SeedSoilMap = ReadRanger(scanner)
	res.SoilFertMap = ReadRanger(scanner)
	res.FertWaterMap = ReadRanger(scanner)
	res.WaterLightMap = ReadRanger(scanner)
	res.LightTempMap = ReadRanger(scanner)
	res.TemphumMap = ReadRanger(scanner)
	res.HumLocMap = ReadRanger(scanner)
	return res
}

func (inp Input) LowestLoc() int {
	locs := []int{}
	for _, seed := range inp.Seeds {
		soil := inp.SeedSoilMap.FindDest(seed)
		fert := inp.SoilFertMap.FindDest(soil)
		water := inp.FertWaterMap.FindDest(fert)
		light := inp.WaterLightMap.FindDest(water)
		temp := inp.LightTempMap.FindDest(light)
		hum:= inp.TemphumMap.FindDest(temp)
		loc := inp.HumLocMap.FindDest(hum)
		locs = append(locs, loc)
	}
	mloc := locs[0]
	for _, v := range locs {
		if v < mloc{ mloc = v}
	}
	return  mloc
}

func Readlines(file io.Reader) (int, int) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	return inp.LowestLoc(), 0
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
