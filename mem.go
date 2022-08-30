package timeseries

import "time"

// memStorage for storing the items
type memStorage[T any, PT Timestamper[T]] struct {
	secStorage
	items []*T
}

// NewMemoryStorage create storage for given types
func NewMemoryStorage[T any, PT Timestamper[T]](period, resolution time.Duration) (Storage[T, PT], error) {
	ms := memStorage[T, PT]{}
	capacity, err := ms.capacityFor(period, resolution)
	if err != nil {
		return nil, err
	}
	ms.items = make([]*T, capacity)

	return &ms, nil
}

// Capacity return number of element can be stored at the storage
func (m *memStorage[T, PT]) Capacity() int {
	return cap(m.items)
}

// At return item at given index
func (m *memStorage[T, PT]) At(idx int64) *T {
	return m.items[idx]
}

// Put set value v at given index idx
func (m *memStorage[T, PT]) Put(idx int64, v *T) error {
	m.items[idx] = v
	return nil
}

// RemoveAt delete items at given storage
func (m *memStorage[T, PT]) RemoveAt(idx int64) error {
	m.items[idx] = nil
	return nil
}

// RemovesOlder removes all items with sequence number less than seq
func (m *memStorage[T, PT]) RemovesOlder(seq int64) int64 {
	n := int64(0)
	for i, v := range m.items {
		if v != nil {
			pv := PT(v)
			if pv.Sequence() < seq {
				m.items[i] = nil
				n++
			}
		}
	}
	return n
}

func (m *memStorage[T, PT]) walkSeq(idx, seq int64, since time.Time, fn func(*T) bool) bool {
	n := int64(len(m.items))
	var v *T
	var pv PT

	// (idx, n)
	for i := idx + 1; i < n; i++ {
		v = m.items[i]
		if v != nil {
			pv = PT(v)
			if pv.Sequence() == seq && pv.Timestamp().After(since) {
				if !fn(v) {
					return false
				}
			}
		}
	}

	// [0, idx]
	for i := int64(0); i <= idx; i++ {
		v = m.items[i]
		if v != nil {
			pv = PT(v)
			if pv.Sequence() == seq && pv.Timestamp().After(since) {
				if !fn(v) {
					return false
				}
			}
		}
	}
	return true
}

// Each iterate over items from oldest to newest
func (m *memStorage[T, PT]) Walk(from int64, fn func(*T) bool) {
	pv := PT(m.items[from])
	// period is inclusive, e.g. if period = 30 minutes, latest value
	// is at 15:30 then all the value since 15:00 will be enumerated.
	// m.period+1 is needed because we are using time.After (> not >=)
	since := pv.Timestamp().Add(-time.Duration(m.period+1) * m.Unit())
	res := m.walkSeq(from, pv.Sequence()-1, since, fn)
	if res {
		m.walkSeq(from, pv.Sequence(), since, fn)
	}
}

func (m *memStorage[T, PT]) eachSeq(idx, seq int64, since time.Time, fn func(*T)) {
	n := int64(len(m.items))
	var v *T
	var pv PT

	// (idx, n)
	for i := idx + 1; i < n; i++ {
		v = m.items[i]
		if v != nil {
			pv = PT(v)
			if pv.Sequence() == seq && pv.Timestamp().After(since) {
				fn(v)
			}
		}
	}

	// [0, idx]
	for i := int64(0); i <= idx; i++ {
		v = m.items[i]
		if v != nil {
			pv = PT(v)
			if pv.Sequence() == seq && pv.Timestamp().After(since) {
				fn(v)
			}
		}
	}
}

// Each iterate over items from oldest to newest
func (m *memStorage[T, PT]) Each(from int64, fn func(*T)) {
	pv := PT(m.items[from])
	// period is inclusive, e.g. if period = 30 minutes, latest value
	// is at 15:30 then all the value since 15:00 will be enumerated.
	// m.period+1 is needed because we are using time.After (> not >=)
	since := pv.Timestamp().Add(-time.Duration(m.period+1) * m.Unit())
	m.eachSeq(from, pv.Sequence()-1, since, fn)
	m.eachSeq(from, pv.Sequence(), since, fn)
}
