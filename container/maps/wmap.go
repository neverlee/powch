package maps

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

type WMap[K comparable, V any] struct {
	data map[K]V
}

func NewWMap[K comparable, V any](capacity ...int) WMap[K, V] {
	c := 16
	if len(capacity) > 0 && capacity[0] > 0 {
		c = capacity[0]
	}

	m := WMap[K, V]{
		data: make(map[K]V, c),
	}

	return m
}

func ToWMap[K comparable, V any](data map[K]V) WMap[K, V] {
	m := WMap[K, V]{
		data: data,
	}

	return m
}

func (m *WMap[K, V]) Set(k K, v V) {
	m.data[k] = v
}

func (m *WMap[K, V]) Get(k K) (V, bool) {
	v, ok := m.data[k]
	return v, ok
}

func (m *WMap[K, V]) Delete(ks ...K) {
	for _, k := range ks {
		delete(m.data, k)
	}
}

func (m *WMap[K, V]) Len() int {
	return len(m.data)
}

func (m *WMap[K, V]) Range(fn func(k K, v V) bool) {
	for k, v := range m.data {
		if !fn(k, v) {
			break
		}
	}
}

func (m *WMap[K, V]) Clone() WMap[K, V] {
	nm := NewWMap[K, V](m.Len())
	for k, v := range m.data {
		nm.data[k] = v
	}
	return nm
}

func (m *WMap[K, V]) Keys() []K {
	ret := make([]K, len(m.data))
	i := 0
	for k := range m.data {
		ret[i] = k
		i++
	}
	return ret
}

func (m *WMap[K, V]) Values() []V {
	ret := make([]V, len(m.data))
	i := 0
	for _, v := range m.data {
		ret[i] = v
		i++
	}
	return ret
}

func (m *WMap[K, V]) Items() []KV[K, V] {
	ret := make([]KV[K, V], len(m.data))
	i := 0
	for k, v := range m.data {
		ret[i].Key = k
		ret[i].Value = v
		i++
	}
	return ret
}

func (m *WMap[K, V]) Raw() map[K]V {
	return m.data
}
