package powch

import (
	"github.com/neverlee/powch/container/list"
)

type KeyFunc[K comparable, T any] func(T) K

type Dedup[K comparable, T any] struct {
	l  *list.List[T]
	kd map[K]*list.Element[T]

	keyFunc KeyFunc[K, T]
	ich     chan T
	och     chan T
	max     int
	closed  bool
	loopEnd bool // nouse, just for test
}

func NewDedup[K comparable, T any](max int, keyFunc KeyFunc[K, T]) *Dedup[K, T] {
	initSize := max
	if initSize < 0 {
		initSize = 0
	}
	dd := &Dedup[K, T]{
		l:  list.New[T](),
		kd: make(map[K]*list.Element[T], max),

		keyFunc: keyFunc,
		ich:     make(chan T, 0),
		och:     make(chan T, 0),
		max:     max,
	}
	go dd.loop()
	return dd
}

func (dd *Dedup[K, T]) loop() {
	var te, to T
	var in, out chan T
	for {
		l := dd.l.Len()
		in = dd.ich
		out = dd.och
		var ele *list.Element[T]
		if l == 0 {
			out = nil
			to = te
		} else {
			ele = dd.l.Front()
			to = ele.Value
		}

		if dd.max > 0 && l >= dd.max {
			in = nil
		}

		if dd.closed {
			in = nil
			if l == 0 {
				close(dd.och)
				break
			}
		}

		select {
		case v, ok := <-in:
			if ok {
				k := dd.keyFunc(v)
				if e, ok := dd.kd[k]; ok {
					e.Value = v
				} else {
					e = dd.l.PushBack(v)
					dd.kd[k] = e
				}
			} else {
				dd.closed = true
			}
		case out <- to:
			if l != 0 {
				k := dd.keyFunc(ele.Value)
				dd.l.Remove(ele)
				delete(dd.kd, k)
			}
		}
	}

	dd.loopEnd = true
}

func (dd *Dedup[K, T]) InChan() chan<- T {
	return dd.ich
}

func (dd *Dedup[K, T]) OutChan() <-chan T {
	return dd.och
}

func (dd *Dedup[K, T]) Close(withClean bool) {
	close(dd.ich)
	if withClean {
		dd.Clean()
	}
}

func (dd *Dedup[K, T]) Clean() {
	for range dd.OutChan() {
	}
}
