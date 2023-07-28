package powch

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Example_PubChannel() {
	// create a pub channel
	pub := NewPub[string]()

	fmt.Println("now start")
	// pub.InChan() <- "will not be showed" // or  pub.Push("will not be showed")

	var wg sync.WaitGroup
	// defer wg.Wait()
	for i := 0; i < 3; i++ {
		wg.Add(1)

		// create a listener
		l := pub.Listen()
		go func(id int, l Listener[string]) {
			defer wg.Done()

			use_chan := id%2 == 0
			if use_chan {
				for {
					if n, ok := <-l.OutChan(); ok {
						d := l.Done(n)
						fmt.Printf("listener %d, get data by chan: {%v}\n", id, d)
					} else {
						fmt.Printf("listener %d done\n", id)
						break
					}
				}
			} else {
				for {
					if d, ok := l.Pop(); ok {
						fmt.Printf("listener %d, get data by pop func: {%v}\n", id, d)
					} else {
						fmt.Printf("listener %d done\n", id)
						break
					}
				}
			}
		}(i, l)
	}

	runtime.Gosched()
	// now you can publish data by this channel
	for i := 0; i < 3; i++ {
		use_chan := i%2 == 0
		fmt.Println("pub", i)
		if use_chan {
			pub.InChan() <- fmt.Sprintf("push a value %d by chan", i)
		} else {
			pub.Push(fmt.Sprintf("push a value %d by push func", i))
		}
	}
	pub.Close()
	wg.Wait()

	// Output:
	// now start
	// pub 0
	// pub 1
	// pub 2
	// listener 2, get data by chan: {push a value 0 by chan}
	// listener 2, get data by chan: {push a value 1 by push func}
	// listener 2, get data by chan: {push a value 2 by chan}
	// listener 2 done
	// listener 1, get data by pop func: {push a value 0 by chan}
	// listener 1, get data by pop func: {push a value 1 by push func}
	// listener 1, get data by pop func: {push a value 2 by chan}
	// listener 1 done
	// listener 0, get data by chan: {push a value 0 by chan}
	// listener 0, get data by chan: {push a value 1 by push func}
	// listener 0, get data by chan: {push a value 2 by chan}
	// listener 0 done
}

func Test_PubChannel(t *testing.T) {
	pub := NewPub[any]()

	l1 := pub.Listen()

	pub.Push("hello")

	{
		d, ok := l1.Pop()
		assert.Equal(t, "hello", d)
		assert.Equal(t, true, ok)
	}

	l2 := pub.Listen()

	pub.Push(123)

	{
		n, ok := <-l1.OutChan()
		d := l1.Done(n)
		assert.Equal(t, 123, d)
		assert.Equal(t, true, ok)
	}

	{
		d, ok := l2.Pop()
		assert.Equal(t, 123, d)
		assert.Equal(t, true, ok)
	}

	pub.Close()

	{
		n, ok := <-l1.OutChan()
		d := l1.Done(n)
		assert.Nil(t, d)
		assert.Equal(t, false, ok)
	}

	{
		_, ok := l2.Pop()
		assert.Equal(t, false, ok)
	}
}
