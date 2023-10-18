package powch

import (
	"sync/atomic"
)

type BNode[T any] struct {
	nextc chan *BNode[T]
	v     T
}

// Publish a publish channel. work with one sender and multi receivers.
// support boardcast datas.
// refer to the design of the linked channel
// see details: https://rogpeppe.wordpress.com/2009/12/01/concurrent-idioms-1-broadcasting-values-in-go-with-linked-channels/)
type Publish[T any] struct {
	curr  atomic.Value
	sendc chan T
}

// NewPublish create a pub channel
func NewPublish[T any]() *Publish[T] {
	bc := Publish[T]{
		sendc: make(chan T),
	}
	bc.curr.Store(make(chan *BNode[T], 1))
	go bc.loop()
	return &bc
}

func (b *Publish[T]) loop() {
	for v := range b.sendc {
		currc := b.curr.Load().(chan *BNode[T])

		nextc := make(chan *BNode[T], 1)
		bn := &BNode[T]{nextc: nextc, v: v}
		currc <- bn
		b.curr.Store(nextc)
	}
	currc := b.curr.Load().(chan *BNode[T])
	close(currc)
}

// Listen get a listener of this pub channel and start listening to the broadcasts
func (b *Publish[T]) Listen() Listener[T] {
	cb := b.curr.Load().(chan *BNode[T])
	return Listener[T]{cb}
}

// InChan returns the send channel.
func (b *Publish[T]) InChan() chan<- T { return b.sendc }

// Push push a data into the send channel
func (b *Publish[T]) Push(v T) { b.sendc <- v }

// Close close this channel
func (b *Publish[T]) Close() {
	close(b.sendc)
}

// Listener a listener
type Listener[T any] struct {
	currc chan *BNode[T]
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

// CaseOut returns the receiver channel that you should only use it in select clause.
// You should always recall CaseOut in each loop
// You should always call CaseDone to finish the operation
func (l *Listener[T]) CaseOut() <-chan *BNode[T] {
	return l.currc
}

// CaseDone unwrap and finish the operation
// if you get a node from CaseOut, even though you don't need the v, you still need to call CaseDone
func (l *Listener[T]) CaseDone(b *BNode[T]) (ret T, ok bool) {
	if b == nil {
		return
	}
	ret, ok = b.v, true
	l.currc <- b
	l.currc = b.nextc
	return
}

// Range calls f sequentially for each value in the chan.
// if f returns false, range stops the iteration.
func (l *Listener[T]) Range(f func(v T) bool) {
	for {
		v, ok := l.CaseDone(<-l.CaseOut())
		if !ok {
			return
		}
		if !f(v) {
			break
		}
	}
}

// returns a new Listener with the same position
func (l *Listener[T]) Clone() Listener[T] {
	return Listener[T]{currc: l.currc}
}
