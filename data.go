package timeseries

import "time"

// Timestamper constraint define Timestamp and pointer to type T
type Timestamper[T any] interface {
	Timestamp() time.Time
	Sequence() int64
	SetSequence(seq int64)
	Merge(other *T) error
	*T
}

// Datum stores information about time series data
type Datum struct {
	Ts  time.Time
	Seq int64
}

// Timestamp return time.Time associated with this data
func (d *Datum) Timestamp() time.Time {
	return d.Ts
}

// Sequence return sequence number for this data
func (d *Datum) Sequence() int64 {
	return d.Seq
}

// SetSequence assign new sequence number for this data
func (d *Datum) SetSequence(seq int64) {
	d.Seq = seq
}
