package powch

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Example_PubChannel() {
	// create a pub channel
	pub := NewPub[string]()

	// create some listener
	l1 := pub.Listen()
	l2 := pub.Listen()

	// now you can publish data by this channel
	pub.InChan() <- "one" // or  pub.Push("one")

	{
		d, ok := l1.Pop()
		fmt.Println("l1 get", d, ok)
	}
	{
		if n, ok := <-l2.OutChan(); ok {
			d := l2.Done(n)
			fmt.Println("l2 get", d, ok)
		}
	}

	// Output:
	// one true
	// one true

	pub.Close()
}

func Test_PubChannel(t *testing.T) {
	pub := NewPub[any]()

	r := pub.Listen()

	pub.Push("hello")

	{
		d, ok := r.Pop()
		assert.Equal(t, "hello", d)
		assert.Equal(t, true, ok)
	}

	r1 := pub.Listen()

	pub.Push(123)

	{
		d, ok := r.Pop()
		assert.Equal(t, 123, d)
		assert.Equal(t, true, ok)
	}

	{
		d, ok := r1.Pop()
		assert.Equal(t, 123, d)
		assert.Equal(t, true, ok)
	}

	pub.Close()

	{
		_, ok := r.Pop()
		assert.Equal(t, false, ok)
	}

	{
		_, ok := r1.Pop()
		assert.Equal(t, false, ok)
	}
}
