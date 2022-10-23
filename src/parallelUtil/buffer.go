package parallelUtil

import "sync"

type Work struct {
	a int
	b int
}

type Buffer struct {
	stop bool
	data []*Work

	mtx  *sync.Mutex
	cond *sync.Cond
}

func NewBuffer() *Buffer {
	b := &Buffer{
		data: make([]*Work, 0),
		stop: false,
		mtx:  &sync.Mutex{},
	}
	b.cond = sync.NewCond(b.mtx)
	return b
}

func (b *Buffer) Push(w *Work) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.data = append(b.data, w)
	b.cond.Broadcast()
}

func (b *Buffer) Pop() *Work {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if len(b.data) <= 0 {
		b.cond.Wait()
	}
	if len(b.data) <= 0 || b.stop {
		return nil
	}
	w := b.data[0]
	b.data = b.data[1:]
	return w
}

func (b *Buffer) Stop() {
	b.stop = true
	b.cond.Broadcast()
}
