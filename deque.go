package powch

import (
	"github.com/neverlee/powch/container/list"
)

// Deque a bouble linkedlist channel
type Deque[T any] struct {
	lst     *list.List[T]
	fin     chan T
	bin     chan T
	fout    chan T
	bout    chan T
	closed  bool
	loopEnd bool
	max     int
}

// NewDeque returns a deque channel
func NewDeque[T any](max int) *Deque[T] {
	dq := &Deque[T]{
		lst:     list.New[T](),
		fin:     make(chan T),
		bin:     make(chan T),
		fout:    make(chan T),
		bout:    make(chan T),
		closed:  false,
		loopEnd: false,
		max:     max,
	}
	go dq.loop()
	return dq
}

func (dq *Deque[T]) loop() {
	var te T
	var tf, tb T
	var fin, fout, bin, bout chan T
	for {
		l := dq.lst.Len()
		fin = dq.fin
		bin = dq.bin
		fout = dq.fout
		bout = dq.bout
		if l == 0 {
			bout, fout = nil, nil
			tf, tb = te, te
		} else {
			tf, tb = dq.lst.Front().Value, dq.lst.Back().Value
		}

		if l == dq.max {
			bin, fin = nil, nil
		}

		if dq.closed {
			bin, fin = nil, nil
			if l == 0 {
				close(dq.fout)
				close(dq.bout)
				break
			}
		}

		select {
		case v, ok := <-fin:
			if ok {
				dq.lst.PushFront(v)
			} else {
				close(dq.bin)
				dq.closed = true
			}
		case v := <-bin:
			dq.lst.PushBack(v)
		case fout <- tf:
			if l != 0 {
				dq.lst.Remove(dq.lst.Front())
			}
		case bout <- tb:
			if l != 0 {
				dq.lst.Remove(dq.lst.Back())
			}
		}
	}

	dq.loopEnd = true
}

func (dq *Deque[T]) FinChan() chan<- T {
	return dq.fin
}

func (dq *Deque[T]) BinChan() chan<- T {
	return dq.bin
}

func (dq *Deque[T]) FoutChan() <-chan T {
	return dq.fout
}

func (dq *Deque[T]) BoutChan() <-chan T {
	return dq.bout
}

func (dq *Deque[T]) Close(withClean bool) {
	close(dq.fin)
	if withClean {
		dq.Clean()
	}
}

func (dq *Deque[T]) Clean() {
	for range dq.FoutChan() {
	}
}
