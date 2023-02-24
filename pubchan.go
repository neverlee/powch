package powch

import (
	"sync/atomic"
)

type BNode[T any] struct {
	nextc chan BNode[T]
	v     T
}

// Pub a pub channel. work with one sender and multi receivers.
// support boardcast datas.
// refer to the design of the linked channel
// see details: https://rogpeppe.wordpress.com/2009/12/01/concurrent-idioms-1-broadcasting-values-in-go-with-linked-channels/)
type Pub[T any] struct {
	curr  atomic.Value
	sendc chan T
}

// NewPub create a pub channel
func NewPub[T any]() *Pub[T] {
	bc := Pub[T]{
		sendc: make(chan T),
	}
	bc.curr.Store(make(chan BNode[T], 1))
	go bc.loop()
	return &bc
}

func (b *Pub[T]) loop() {
	for v := range b.sendc {
		currc := b.curr.Load().(chan BNode[T])

		nextc := make(chan BNode[T], 1)
		bn := BNode[T]{nextc: nextc, v: v}
		currc <- bn
		b.curr.Store(nextc)
	}
	currc := b.curr.Load().(chan BNode[T])
	close(currc)
}

// Listen get a listener of this pub channel and start listening to the broadcasts
func (b *Pub[T]) Listen() Listener[T] {
	cb := b.curr.Load().(chan BNode[T])
	return Listener[T]{cb}
}

// InChan returns the send channel.
func (b *Pub[T]) InChan() chan<- T { return b.sendc }

// Push push a data into the send channel
func (b *Pub[T]) Push(v T) { b.sendc <- v }

// Close close this channel
func (b *Pub[T]) Close() {
	close(b.sendc)
}

// Listener a listener
type Listener[T any] struct {
	currc chan BNode[T]
}

// Pop get value from the pub channel.
// returns a value and true, if there are still datas in the pub channel.
// returns a empty value and false, if there is no datas in the pub channel and the channal is closed
// otherwise block
func (l *Listener[T]) Pop() (T, bool) {
	b, ok := <-l.currc
	if !ok {
		var v T
		return v, ok
	}

	v := b.v
	l.currc <- b
	l.currc = b.nextc
	return v, ok
}

// OutChan returns the receiver channel that you can get data from it.
// You shoud always call Done to finish the operation
func (l *Listener[T]) OutChan() <-chan BNode[T] {
	return l.currc
}

// Done unwrap and finish the operation
// if you use OutChan, you should always use this func
func (l *Listener[T]) Done(b BNode[T]) T {
	v := b.v
	if b.nextc == nil {
		return v
	}

	l.currc <- b
	l.currc = b.nextc
	return v
}
