package serialUtil

import (
	"github.com/jaydenhthompson/bst_compare/src/bstUtil"
	"github.com/jaydenhthompson/bst_compare/src/tree"
)

func ClassifyTreeHashes(t []*tree.Tree, m map[int][]int) {
	for i, tree := range t {
		hash := tree.CalculateHash()
		m[hash] = append(m[hash], i)
	}
}

func GroupHashedTrees(trees []*tree.Tree, hashMap map[int][]int) [][]int {
	groups := make([][]int, 0)
	for _, matchingTrees := range hashMap {
		groups = append(groups, compareTreeArray(trees, matchingTrees)...)
	}
	return groups
}

func compareTreeArray(trees []*tree.Tree, indexes []int) [][]int {
	matchMap := make(map[int]map[int]struct{})
	matched := make([]bool, len(indexes))
	for i := range indexes {
		for j := i + 1; j < len(indexes); j++ {
			if matched[j] {
				continue
			}
			if bstUtil.CompareSlices(trees[indexes[i]].InOrderTraversal(), trees[indexes[j]].InOrderTraversal()) {
				if matchMap[indexes[i]] == nil {
					matchMap[indexes[i]] = map[int]struct{}{
						indexes[i]: {},
					}
				}
				matched[j] = true
				matchMap[indexes[i]][indexes[j]] = struct{}{}
			}
		}
	}

	matchMatrix := make([][]int, len(matchMap))
	i := 0
	for _, v := range matchMap {
		matchMatrix[i] = make([]int, 0, len(v))
		for k, _ := range v {
			matchMatrix[i] = append(matchMatrix[i], k)
		}
		i++
	}
	return matchMatrix
}
