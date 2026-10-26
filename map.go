package toolkit

import "sync"

// 安全map再封装

type Map[K comparable, V any] struct {
	data sync.Map
}

func (m *Map[K, V]) Set(key K, value V) {
	m.data.Store(key, value)
}

func (m *Map[K, V]) MSet(kv map[K]V) {
	for k, v := range kv {
		m.Set(k, v)
	}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.data.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

func (m *Map[K, V]) MGet(keys []K) map[K]V {
	r := make(map[K]V, len(keys))
	for _, k := range keys {
		v, ok := m.data.Load(k)
		if ok {
			r[k] = v.(V)
		}
	}

	return r
}

func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.data.Load(key)
	return ok
}

func (m *Map[K, V]) Del(key K) {
	m.data.Delete(key)
}
