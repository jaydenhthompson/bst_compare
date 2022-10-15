package bstUtil

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/jaydenhthompson/bst_compare/src/tree"
)

func ParseFile(inputFile string) []*tree.Tree {
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

func CompareSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func PruneMap(m map[int][]int) {
	for key, arr := range m {
		if len(arr) <= 1 {
			delete(m, key)
		}
	}
}
