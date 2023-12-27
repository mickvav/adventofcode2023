package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Input struct {
	m       []string
	nodes map[string]*node
}

type node struct {
	name string
	edges map[string]*edge
}

type edge struct {
	n1,n2 *node
	weight int
}


func (inp *Input) parseGraph() {
	inp.nodes = map[string]*node{}
	for _, line := range inp.m {
		p:= strings.Split(line, ": ")
		v:= p[0]
		if _, ok := inp.nodes[v]; !ok {
			inp.nodes[v] = &node{
				name:  v,
				edges: map[string]*edge{},
			}
		}
		for _, v2 := range strings.Split(p[1], " ") {
			if _, ok:=inp.nodes[v2]; !ok {
				inp.nodes[v2] = &node{
					name:  v2,
					edges: map[string]*edge{},
				}
			}
			e := edge{
				n1:     inp.nodes[v],
				n2:     inp.nodes[v2],
				weight: 1,
			}
			inp.nodes[v].edges[v2] = &e
			inp.nodes[v2].edges[v] = &e
		}	
	}
}

func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m:     []string{},
		nodes: map[string]*node{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)

	}
	res.parseGraph()
	return res
}

func (e *edge) other(n *node) *node {
	if e.n1 == n {
		return e.n2
	} else {
		return e.n1
	}
}
func (inp *Input) Shrinknodes(n1, n2 string) error {
	if n1obj, ok := inp.nodes[n1]; !ok {
		return fmt.Errorf("not found %s", n1)
	} else if n2obj, ok := inp.nodes[n2]; !ok {
		return fmt.Errorf("not found %s", n2)
	} else if _, ok:= n1obj.edges[n2]; !ok{
		return fmt.Errorf("not neighbours: %s %s", n1, n2)
	} else {
		newnodename := n1+n2
		newnode := node{
			name:  newnodename,
			edges: map[string]*edge{},
		}
		for k, e2 := range n1obj.edges {
			if k != n2 {
				nremote := e2.other(n1obj)
				enew := edge{
					n1:     &newnode,
					n2:     nremote,
					weight: e2.weight,
				}
				newnode.edges[k] = &enew
				nremote.edges[newnodename] = &enew
				delete(nremote.edges, n1)
			}
		}
		for k, e2 := range n2obj.edges {
			if k != n1 {
				nremote := e2.other(n2obj)
				var enew *edge
				if _, ok := newnode.edges[k]; !ok {
					enew = &edge{
						n1: &newnode,
						n2: nremote,
						weight: e2.weight,
					}
					newnode.edges[k] = enew
				} else {
					enew = newnode.edges[k]
					enew.weight += e2.weight
				}
				nremote.edges[newnodename] = enew
				delete(nremote.edges, n2)
			}
		}
		inp.nodes[newnodename] = &newnode
		for k := range n1obj.edges {
			delete(n1obj.edges, k)
		}
		for k := range n2obj.edges {
			delete(n2obj.edges, k)
		}
		delete(inp.nodes, n1)
		delete(inp.nodes, n2)
		return nil
	}
}

func (inp *Input) ShrinkToLimit() {
	for len(inp.nodes) > 2 {
		keys := make([]string, 0, len(inp.nodes))
		for k := range inp.nodes {
			keys = append(keys, k)
		}
		idx:=rand.Int() % len(keys)
		keys2 := make([]string, 0, len(inp.nodes[keys[idx]].edges))
		for k := range inp.nodes[keys[idx]].edges {
			keys2 = append(keys2, k)
		}
		idx2 := rand.Int() % len(keys2)
		inp.Shrinknodes(keys[idx], keys2[idx2])
	}
}
func (inp *Input) Print() {
	for _, v := range inp.m {
		fmt.Printf(string(v) + "\n")
	}
}

func (inp *Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	for i:=1; i<len(inp.m)*len(inp.m); i++ {
		inp.ShrinkToLimit()
		keys := make([]string, 0, len(inp.nodes))
		for k := range inp.nodes {
			keys = append(keys, k)
		}
		flux := 0
		for _, e := range inp.nodes[keys[0]].edges {
			flux += e.weight
		}
		if flux == 3 {
			fmt.Printf("found! %s %s\n", keys[0], keys[1])
			s1 = (len(keys[0])*len(keys[1]))/9
			break
		}

		inp.parseGraph()
	}
	return s1,s2
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
