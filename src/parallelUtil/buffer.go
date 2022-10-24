package parallelUtil

import (
	"container/list"
	"sync"
)

type Work struct {
	a int
	b int
}

type Buffer struct {
	stop bool
	data *list.List

	mtx     *sync.Mutex
	popCond *sync.Cond
}

func NewBuffer() *Buffer {
	b := &Buffer{
		data: list.New(),
		stop: false,
		mtx:  &sync.Mutex{},
	}
	b.popCond = sync.NewCond(b.mtx)
	return b
}

func (b *Buffer) Push(w *Work) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.data.PushBack(w)
	b.popCond.Signal()
}

func (b *Buffer) Pop() *Work {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	for b.data.Len() <= 0 {
		if b.stop {
			return nil
		}
		b.popCond.Wait()
	}
	cur := b.data.Front()
	w := cur.Value.(*Work)
	b.data.Remove(cur)
	return w
}

func (b *Buffer) Wait() {
	for {
		if b.data.Len() <= 0 {
			return
		}
	}
}

func (b *Buffer) Stop() {
	b.stop = true
	b.popCond.Broadcast()
}
