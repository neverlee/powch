package powch

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Example_Publish() {
	// create a pub channel
	pub := NewPublish[string]()

	fmt.Println("now start")

	pub.InChan() <- "will not be showed" // or  pub.Push("will not be showed")

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(3)

	// create a listener
	l := pub.Listen()
	{
		go func(id int, l Listener[string]) {
			defer wg.Done()
			for {
				select {
				case n := <-l.CaseOut():
					if d, ok := l.CaseDone(n); ok {
						fmt.Printf("listener %d, get data by chan: {%v}\n", id, d)
					} else {
						fmt.Printf("listener %d done\n", id)
						return
					}
				}
			}
		}(0, l.Clone())
	}

	{
		go func(id int, l Listener[string]) {
			defer wg.Done()
			for {
				if d, ok := l.Pop(); ok {
					fmt.Printf("listener %d, get data by pop func: {%v}\n", id, d)
				} else {
					fmt.Printf("listener %d done\n", id)
					break
				}
			}
		}(1, l.Clone())
	}

	{
		go func(id int, l Listener[string]) {
			defer wg.Done()
			fmt.Printf("listener %d done\n", id)
			l.Range(func(d string) bool {
				fmt.Printf("listener %d, get data by pop func: {%v}\n", id, d)
				return true
			})
			fmt.Printf("listener %d done\n", id)
		}(2, l.Clone())
	}

	// You should not use it like these
	/*
		{
			l := pub.Listen()

			// one: use it in for range
			for bn := range l.CaseOut() {
				v := l.CaseDone(bn)
				// ...
			}

			// two: use the chan var in a loop
			ch := l.CaseOut()
			for {
				bn := <-ch
				v := l.CaseDone(bn)
				// ...
			}

			// three: without CaseDone
			for {
				<-l.CaseOut()
				// ...
			}
		}
	*/

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
}

func Test_Publish(t *testing.T) {
	pub := NewPublish[any]()

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
		d, ok := l1.CaseDone(<-l1.CaseOut())
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
		d, ok := l1.CaseDone(<-l1.CaseOut())
		assert.Nil(t, d)
		assert.Equal(t, false, ok)
	}

	{
		_, ok := l2.Pop()
		assert.Equal(t, false, ok)
	}
}
