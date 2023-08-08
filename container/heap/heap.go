package heap

type Heap[T any] struct {
	data []T
	cmp  func(a, b T) bool
}

func NewHeap[T any](sz int, less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{
		data: make([]T, 0, sz),
		cmp:  less,
	}
	return h
}

func (h *Heap[T]) Len() int {
	return len(h.data)
}

func (h *Heap[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *Heap[T]) Less(i, j int) bool {
	return h.cmp(h.data[i], h.data[j])
}

func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	h.up(h.Len() - 1)
}

func (h *Heap[T]) Top() T {
	return h.data[0]
}

func (h *Heap[T]) Pop() T {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)

	r := h.data[n]
	h.data = h.data[:n]
	return r
}

func (h *Heap[T]) Init() {
	n := len(h.data)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (h *Heap[T]) Remove(i int) any {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.Pop()
}

func (h *Heap[T]) Fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

func (h *Heap[T]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func (h *Heap[T]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}
