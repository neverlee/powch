package powch



var Lock = Unit{}

type Mutex struct {
	ch chan Unit
}

func NewMutex() Mutex {
	return Mutex{
		ch: make(chan Unit, 1),
	}
}
