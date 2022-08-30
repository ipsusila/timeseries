package timeseries

import (
	"math"
	"time"
)

// nullStorage for storing the items
type nullStorage[T any, PT Timestamper[T]] struct {
	secStorage
}

// NewNullStorage create storage for given types
func NewNullStorage[T any, PT Timestamper[T]](period, resolution time.Duration) (Storage[T, PT], error) {
	ns := &nullStorage[T, PT]{}
	_, err := ns.capacityFor(period, resolution)
	if err != nil {
		return nil, err
	}

	return ns, nil
}

// Capacity return number of element can be stored at the storage
func (s *nullStorage[T, PT]) Capacity() int {
	return math.MaxInt
}

// At return item at given index
func (s *nullStorage[T, PT]) At(idx int64) *T {
	return nil
}

// Put set value v at given index idx
func (s *nullStorage[T, PT]) Put(idx int64, v *T) error {
	return nil
}

// RemoveAt delete items at given storage
func (s *nullStorage[T, PT]) RemoveAt(idx int64) error {
	return nil
}
func (s *nullStorage[T, PT]) RemovesOlder(seq int64) int64 {
	return 0
}
func (s *nullStorage[T, PT]) Each(from int64, fn func(*T)) {
	// do nothing
}
func (s *nullStorage[T, PT]) Walk(from int64, fn func(*T) bool) {
	// do nothing
}
