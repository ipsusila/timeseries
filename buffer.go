package timeseries

// Buffer for storing data point
type Buffer[T any, PT Timestamper[T]] interface {
	Capacity() int
	Put(v *T) error
	Latest() *T
	AttachStorage(s Storage[T, PT])
	Vacuum() error
	Each(fn func(*T))
	Walk(fn func(*T) bool)
}
