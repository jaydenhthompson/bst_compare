package parallelUtil

import (
	"sync"

	"github.com/jaydenhthompson/bst_compare/src/bstUtil"
	"github.com/jaydenhthompson/bst_compare/src/tree"
)

var (
	wg sync.WaitGroup

	dataMut  sync.Mutex
	hashStop chan bool
	dataStop chan bool

	running bool = true
	mtx          = &sync.Mutex{}
	cond         = sync.NewCond(mtx)
)

type Data struct {
	idx  int
	hash int
}

func DataWorker(m map[int][]int, dataChan chan Data) {
	for {
		select {
		case data := <-dataChan:
			dataMut.Lock()
			m[data.hash] = append(m[data.hash], data.idx)
			dataMut.Unlock()
		case <-dataStop:
			wg.Done()
			return
		}
	}
}

func HashWorker(t []*tree.Tree, index chan int, dataChan chan Data) {
	for {
		select {
		case idx := <-index:
			hash := t[idx].CalculateHash()
			dataChan <- Data{idx: idx, hash: hash}
		case <-hashStop:
			wg.Done()
			return
		}
	}
}

func FirstHashImplementation(t []*tree.Tree, m map[int][]int, dataWorkers int) {
	dataStop = make(chan bool)
	dataChan := make(chan Data)

	wg.Add(len(t))
	for i := 0; i < dataWorkers; i++ {
		go DataWorker(m, dataChan)
	}
	for i := 0; i < len(t); i++ {
		go func(i int) {
			defer wg.Done()
			hash := t[i].CalculateHash()
			dataChan <- Data{idx: i, hash: hash}
		}(i)
	}
	wg.Wait()

	wg.Add(dataWorkers)
	for i := 0; i < dataWorkers; i++ {
		dataStop <- true
	}
	wg.Wait()
}

func ClassifyTreeHashes(t []*tree.Tree, m map[int][]int, hashWorkers, dataWorkers int) {
	dataStop = make(chan bool)
	hashStop = make(chan bool)
	dataChan := make(chan Data)
	indexChan := make(chan int)

	wg.Add(hashWorkers)
	for i := 0; i < hashWorkers; i++ {
		go HashWorker(t, indexChan, dataChan)
	}
	for i := 0; i < dataWorkers; i++ {
		go DataWorker(m, dataChan)
	}

	for i := range t {
		indexChan <- i
	}
	for i := 0; i < hashWorkers; i++ {
		hashStop <- true
	}
	wg.Wait()

	wg.Add(dataWorkers)
	for i := 0; i < dataWorkers; i++ {
		dataStop <- true
	}
	wg.Wait()
}

func CompWorker(trees []*tree.Tree, buff *Buffer, adj [][]bool, matched []bool) {
	for {
		w := buff.Pop()
		if w != nil {
			if adj[w.b][w.a] {
				continue
			}
			eq := bstUtil.CompareSlices(trees[w.a].InOrderTraversal(), trees[w.b].InOrderTraversal())
			if eq {
				adj[w.a][w.b] = true
				adj[w.b][w.a] = true
			}
		} else if !running {
			break
		}
	}
	wg.Done()
}

func FirstCompWorker(trees []*tree.Tree, adj [][]bool, i, j int) {
	eq := bstUtil.CompareSlices(trees[i].InOrderTraversal(), trees[j].InOrderTraversal())
	if eq {
		adj[i][j] = true
		adj[j][i] = true
	}
}

func GroupHashedTrees(trees []*tree.Tree, hashMap map[int][]int, compWorkers int) [][]int {
	adj := make([][]bool, len(trees))
	matched := make([]bool, len(trees))
	for i := range adj {
		adj[i] = make([]bool, len(trees))
	}

	wg.Add(compWorkers)
	buff := NewBuffer()
	for i := 0; i < compWorkers; i++ {
		go CompWorker(trees, buff, adj, matched)
	}

	for _, matchingTrees := range hashMap {
		for i := range matchingTrees {
			for j := i + 1; j < len(matchingTrees); j++ {
				buff.Push(&Work{a: matchingTrees[i], b: matchingTrees[j]})
			}
		}
	}

	buff.Wait()
	running = false
	buff.Stop()
	wg.Wait()

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
