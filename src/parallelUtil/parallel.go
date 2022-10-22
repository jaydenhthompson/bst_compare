package parallelUtil

import (
	"sync"

	"github.com/jaydenhthompson/bst_compare/src/tree"
)

var (
	wg sync.WaitGroup

	dataMut  sync.Mutex
	hashStop chan bool
	dataStop chan bool
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
