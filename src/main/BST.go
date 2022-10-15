package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jaydenhthompson/bst_compare/src/bstUtil"
	"github.com/jaydenhthompson/bst_compare/src/tree"
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

func classifyTreeHashes(t []*tree.Tree, m map[int][]int) {
	for i, tree := range t {
		hash := tree.CalculateHash()
		m[hash] = append(m[hash], i)
	}
}

func compareHashedTrees(trees []*tree.Tree, hashMap map[int][]int, groupMap map[int][]int) {

}

func printHashMapping(time time.Duration, m map[int][]int) {
	fmt.Printf("hashGroupTime: %d\n", time.Microseconds())
	i := 0
	for _, arr := range m {
		fmt.Printf("hash%d:", i)
		for _, index := range arr {
			fmt.Printf(" id%d", index)
		}
		println()
		i++
	}
}

func main() {
	trees := bstUtil.ParseFile(inputFile)
	hashMap := make(map[int][]int)

	hashStart := time.Now()
	classifyTreeHashes(trees, hashMap)
	hashDuration := time.Since(hashStart)

	// util for removing hashes with only one tree
	bstUtil.PruneMap(hashMap)

	//groupMap := make(map[int][]int)

	//compareStart := time.Now()
	//compareHashedTrees(trees, hashMap, groupMap)
	//compareDuration := time.Since(compareStart)
	//bstUtil.PruneMap(groupMap)

	printHashMapping(hashDuration, hashMap)
}
