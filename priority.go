package powch

import "github.com/neverlee/powch/container/heap"

type PriorityQueue[T any] struct {
	h       *heap.Heap[T]
	ich     chan T
	och     chan T
	max     int
	closed  bool
	loopEnd bool
}

func NewPriorityQueue[T any](max int, less func(a, b T) bool) *PriorityQueue[T] {
	initSize := max
	if initSize < 0 {
		initSize = 0
	}

	pq := &PriorityQueue[T]{
		h:   heap.NewHeap[T](initSize, less),
		ich: make(chan T, 0),
		och: make(chan T, 0),
		max: max,
	}
	go pq.loop()
	return pq
}

func (pq *PriorityQueue[T]) loop() {
	var te, to T
	var in, out chan T
	for {
		l := pq.h.Len()
		in = pq.ich
		out = pq.och
		if l == 0 {
			out = nil
			to = te
		} else {
			to = pq.h.Top()
		}

		if pq.max > 0 && l >= pq.max {
			in = nil
		}

		if pq.closed {
			in = nil
			if l == 0 {
				close(pq.och)
				break
			}
		}

		select {
		case v, ok := <-in:
			if ok {
				pq.h.Push(v)
			} else {
				pq.closed = true
			}
		case out <- to:
			if l != 0 {
				pq.h.Pop()
			}
		}
	}

	pq.loopEnd = true
}

func (pq *PriorityQueue[T]) InChan() chan<- T {
	return pq.ich
}

func (pq *PriorityQueue[T]) OutChan() <-chan T {
	return pq.och
}

func (pq *PriorityQueue[T]) Close(withClean bool) {
	close(pq.ich)
	if withClean {
		pq.Clean()
	}
}

func (pq *PriorityQueue[T]) Clean() {
	for range pq.OutChan() {
	}
}
