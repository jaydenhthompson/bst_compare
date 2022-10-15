package main

import (
	"flag"
	"io/ioutil"
	"strconv"
	"strings"

	"src/main/src/tree"
)

var (
	inputFile   string
	hashWorkers int
	dataWorkers int
	compWorkers int
)

func init() {
	flag.StringVar(&inputFile, "input", "", "input file")
	flag.IntVar(&hashWorkers, "hash-workers", 1, "number of hash workers")
	flag.IntVar(&dataWorkers, "data-workers", 1, "number of data workers")
	flag.IntVar(&compWorkers, "comp-workers", 1, "number of comp workers")

	flag.Parse()
}

func parseFile() []*tree.Tree {
	trees := []*tree.Tree{}
	body, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		bst := &tree.Tree{}
		values := strings.Split(line, " ")
		for _, value := range values {
			if value == "" {
				continue
			}
			intVal, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			bst.Add(intVal)
		}
		trees = append(trees, bst)
	}

	return trees
}

func main() {
	trees := parseFile()
	for _, tree := range trees {
		for _, val := range tree.InOrderTraversal() {
			print(val)
			print(" ")
		}
		println()
	}

	for _, tree := range trees {
		println(tree.CalculateHash())
	}
}
