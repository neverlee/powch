package maps

import (
	"github.com/neverlee/powch/container/mutex"
)

type SyncWMap[K comparable, V any] struct {
	data WMap[K, V]
	mu   mutex.ExMutex
}

func NewSyncWMap[K comparable, V any](capacity ...int) *SyncWMap[K, V] {
	m := SyncWMap[K, V]{
		data: NewWMap[K, V](capacity...),
	}
	return &m
}

func ToSyncWMap[K comparable, V any](data WMap[K, V]) *SyncWMap[K, V] {
	m := SyncWMap[K, V]{
		data: data,
	}
	return &m
}

func (m *SyncWMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data.Set(k, v)
}

func (m *SyncWMap[K, V]) Get(k K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data.Get(k)
}

func (m *SyncWMap[K, V]) Delete(ks ...K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data.Delete(ks...)
}

func (m *SyncWMap[K, V]) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.data.data)
}

func (m *SyncWMap[K, V]) Range(fn func(k K, v V) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data.Range(fn)
}

func (m *SyncWMap[K, V]) Clone() *SyncWMap[K, V] {
	m.mu.Lock()
	defer m.mu.Unlock()

	return ToSyncWMap(m.data.Clone())
}

func (m *SyncWMap[K, V]) Keys() []K {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.data.Keys()
}

func (m *SyncWMap[K, V]) Values() []V {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.data.Values()
}

func (m *SyncWMap[K, V]) Items() []KV[K, V] {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.data.Items()
}

func (m *SyncWMap[K, V]) Eject() WMap[K, V] {
	m.mu.Lock()
	defer m.mu.Unlock()

	om := m.data
	m.data = NewWMap[K, V]()

	return om
}

func (m *SyncWMap[K, V]) Do(fn func()) {
	m.mu.Apply(fn)
}
