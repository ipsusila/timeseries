package timeseries

import "time"

// StorageType
type StorageType int

// Type of storages
const (
	NullStore StorageType = iota
	MemStore
	FileStore
)

// Storage for storing items
type Storage[T any, PT Timestamper[T]] interface {
	Capacity() int
	At(idx int64) *T
	Put(idx int64, v *T) error
	RemoveAt(idx int64) error
	RemovesOlder(seq int64) int64
	Period() int64
	Resolution() int64
	Unit() time.Duration
	Each(from int64, fn func(*T))
	Walk(from int64, fn func(*T) bool)
}

// secStorage stores time series data in seconds resolution
type secStorage struct {
	period     int64
	resolution int64
}

func (b *secStorage) Unit() time.Duration {
	return time.Second
}

func (b *secStorage) Period() int64 {
	return b.period
}
func (b *secStorage) Resolution() int64 {
	return b.resolution
}
func (b *secStorage) capacityFor(period, resolution time.Duration) (int64, error) {
	if resolution < time.Second {
		return -1, ErrInvalidResolution
	}
	if resolution >= period {
		return -1, ErrInvalidPeriod
	}

	secPeriod := int64(period.Seconds())
	secResolution := int64(resolution.Seconds())
	capacity := secPeriod / secResolution
	if secPeriod%secResolution != 0 {
		capacity++
		secPeriod = capacity * secResolution
	}

	b.period = secPeriod
	b.resolution = secResolution
	return capacity, nil
}
