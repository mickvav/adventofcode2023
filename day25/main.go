package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/daviddengcn/go-algs/maxflow"
)

type Input struct {
	m     []string
	graph maxflow.Graph
	nodes map[string]*maxflow.Node
}

func ParseLine(l string) (string, []string) {
	p := strings.Split(l, ": ")
	p1 := strings.Split(p[1], " ")
	return p[0], p1
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:     []string{},
		graph: *maxflow.NewGraph(),
		nodes: map[string]*maxflow.Node{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		node, connected := ParseLine(line)
		if _, ok := res.nodes[node]; !ok {
			res.nodes[node] = res.graph.AddNode()
		}
		for _, v := range connected {
			if _, ok := res.nodes[v]; !ok {
				res.nodes[v] = res.graph.AddNode()
			}
			res.graph.AddEdge(res.nodes[node], res.nodes[v], 1, 2)
		}
	}
	for _, n := range res.nodes {
		res.graph.SetTweights(n, 1, 2)
	}
	return res
}

func (inp *Input) GenGraph(n1, n2 string) {
	inp.graph = *maxflow.NewGraph()
	inp.nodes = make(map[string]*maxflow.Node)
	for _, line := range inp.m {
		node, connected := ParseLine(line)
		if _, ok := inp.nodes[node]; !ok {
			inp.nodes[node] = inp.graph.AddNode()
		}
		for _, v := range connected {
			if _, ok := inp.nodes[v]; !ok {
				inp.nodes[v] = inp.graph.AddNode()
			}
			inp.graph.AddEdge(inp.nodes[node], inp.nodes[v], 1, 1)
			inp.graph.AddEdge(inp.nodes[v], inp.nodes[node], 1, 1)
		}
	}
	src := inp.graph.AddNode()
	inp.graph.AddEdge(src, inp.nodes[n1], math.MaxInt32, 0)
	dst := inp.graph.AddNode()
	inp.graph.AddEdge(inp.nodes[n2], dst, 0, math.MaxInt32)

	inp.graph.SetTweights(inp.nodes[n1], 1, 0)
	inp.graph.SetTweights(inp.nodes[n2], 0, 1)
}

func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	keys := make([]string, 0, len(inp.nodes))
	for k := range inp.nodes {
		keys = append(keys, k)
	}
	for i1, v1 := range keys {
		for i2, v2 := range keys {
			if i1 != i2 {
				inp.GenGraph(v1, v2)
				inp.graph.Run()
				fmt.Printf("%s %s %d\n", v1, v2, inp.graph.Flow())
				sources := 0
				sinks := 0
				for _, n := range inp.nodes {
					if inp.graph.IsSource(n) {
						sources += 1
					} else {
						sinks += 1
					}
				}
				fmt.Printf("  %d %d %d\n", sources, sinks, sources*sinks)
			}
		}
	}
	sources := 0
	sinks := 0
	for _, n := range inp.nodes {
		if inp.graph.IsSource(n) {
			sources += 1
		} else {
			sinks += 1
		}
	}
	fmt.Printf("%d\n", inp.graph.Flow())
	s1 = sources * sinks
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
