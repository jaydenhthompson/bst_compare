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
	adj := make([][]bool, len(trees))
	for i := range adj {
		adj[i] = make([]bool, len(trees))
	}

	for _, matchingTrees := range hashMap {
		for i := range matchingTrees {
			for j := i + 1; j < len(matchingTrees); j++ {
				if adj[matchingTrees[j]][matchingTrees[i]] {
					continue
				}
				if bstUtil.CompareSlices(trees[matchingTrees[i]].InOrderTraversal(), trees[matchingTrees[j]].InOrderTraversal()) {
					adj[matchingTrees[i]][matchingTrees[j]] = true
					adj[matchingTrees[j]][matchingTrees[i]] = true
				}
			}
		}
	}

	visited := make(map[int]bool)
	groups := make([][]int, 0)
	for i := 0; i < len(adj); i++ {
		if visited[i] {
			continue
		}
		group := make([]int, 0)
		group = append(group, i)
		for j := 0; j < len(adj); j++ {
			if !adj[i][j] || visited[j] {
				continue
			}
			visited[j] = true
			group = append(group, j)
		}
		if len(group) > 1 {
			groups = append(groups, group)
		}
		visited[i] = true
	}
	return groups
}
