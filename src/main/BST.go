package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jaydenhthompson/bst_compare/src/bstUtil"
	"github.com/jaydenhthompson/bst_compare/src/parallelUtil"
	"github.com/jaydenhthompson/bst_compare/src/serialUtil"
)

var (
	inputFile   string
	hashWorkers int
	dataWorkers int
	compWorkers int
	timeOnly    bool
)

func init() {
	flag.StringVar(&inputFile, "input", "", "input file")
	flag.IntVar(&hashWorkers, "hash-workers", 0, "number of hash workers")
	flag.IntVar(&dataWorkers, "data-workers", 1, "number of data workers")
	flag.IntVar(&compWorkers, "comp-workers", 0, "number of comp workers")
	flag.BoolVar(&timeOnly, "t", false, "print times only")
	flag.Parse()
}

func printHashMapping(time time.Duration, m map[int][]int) {
	fmt.Printf("hashGroupTime: %d\n", time.Microseconds())
	if timeOnly {
		return
	}
	i := 0
	for _, arr := range m {
		fmt.Printf("hash %d:", i)
		for _, index := range arr {
			fmt.Printf(" %d", index)
		}
		fmt.Println()
		i++
	}
}

func printGroupings(time time.Duration, groups [][]int) {
	fmt.Printf("compareTreeTime: %d\n", time.Microseconds())
	if timeOnly {
		return
	}
	for i, group := range groups {
		fmt.Printf("group %d:", i)
		for _, id := range group {
			fmt.Printf(" %d", id)
		}
		fmt.Println()
	}
}

func main() {
	trees := bstUtil.ParseFile(inputFile)
	hashMap := make(map[int][]int)

	hashStart := time.Now()
	if hashWorkers <= 1 {
		serialUtil.ClassifyTreeHashes(trees, hashMap)
	} else {
		parallelUtil.ClassifyTreeHashes(trees, hashMap, hashWorkers, dataWorkers)
	}
	hashDuration := time.Since(hashStart)
	bstUtil.PruneMap(hashMap)
	printHashMapping(hashDuration, hashMap)

	if dataWorkers > 0 && compWorkers > 0 {
		var groups [][]int
		compareStart := time.Now()
		if compWorkers <= 1 {
			groups = serialUtil.GroupHashedTrees(trees, hashMap)
		} else {
			groups = parallelUtil.GroupHashedTrees(trees, hashMap, compWorkers)
		}
		compareDuration := time.Since(compareStart)
		printGroupings(compareDuration, groups)
	}
}
