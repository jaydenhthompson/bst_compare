package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jaydenhthompson/bst_compare/src/bstUtil"
	"github.com/jaydenhthompson/bst_compare/src/serialUtil"
)

var (
	inputFile   string
	hashWorkers int
	dataWorkers int
	compWorkers int
)

func init() {
	flag.StringVar(&inputFile, "input", "", "input file")
	flag.IntVar(&hashWorkers, "hash-workers", 0, "number of hash workers")
	flag.IntVar(&dataWorkers, "data-workers", 0, "number of data workers")
	flag.IntVar(&compWorkers, "comp-workers", 0, "number of comp workers")

	flag.Parse()
}

func printHashMapping(time time.Duration, m map[int][]int) {
	fmt.Printf("hashGroupTime: %d\n", time.Microseconds())
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
	serialUtil.ClassifyTreeHashes(trees, hashMap)
	hashDuration := time.Since(hashStart)

	// util for removing hashes with only one tree
	bstUtil.PruneMap(hashMap)

	compareStart := time.Now()
	groups := serialUtil.GroupHashedTrees(trees, hashMap)
	compareDuration := time.Since(compareStart)

	printHashMapping(hashDuration, hashMap)
	printGroupings(compareDuration, groups)
}
