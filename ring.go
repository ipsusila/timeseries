package timeseries

import (
	"errors"
	"time"
)

// List of known errors
var (
	ErrInvalidResolution = errors.New("invalid resolution")
	ErrInvalidPeriod     = errors.New("invalid period")
	ErrOutOfRange        = errors.New("timestamp out of range")
	ErrUnknownStorage    = errors.New("unknown storage type")
)

const (
	vInvalid = -1
)

type ringBuf[T any, PT Timestamper[T]] struct {
	maxEpoch int64
	maxIdx   int64
	latest   *T
	items    Storage[T, PT]
}

// NewBuffer create circular buffer with given time resolution and period of time.
func NewBuffer[T any, PT Timestamper[T]](period, resolution time.Duration, st StorageType) (Buffer[T, PT], error) {
	// Allocate storage
	var err error
	var stg Storage[T, PT]
	switch st {
	case NullStore:
		stg, err = NewNullStorage[T, PT](period, resolution)
	case MemStore:
		stg, err = NewMemoryStorage[T, PT](period, resolution)
	default:
		err = ErrUnknownStorage
	}
	if err != nil {
		return nil, err
	}

	return &ringBuf[T, PT]{
		maxEpoch: vInvalid,
		maxIdx:   vInvalid,
		latest:   nil,
		items:    stg,
	}, nil
}

// calculate position from epoch
func (r *ringBuf[T, PT]) pos(epoch int64) (int64, int64) {
	seq := epoch / r.items.Period()
	idx := (epoch % r.items.Period()) / r.items.Resolution()

	return seq, idx
}

func (r *ringBuf[T, PT]) Latest() *T {
	return r.latest
}

// Capacity return capacity for this ring buffer
func (r *ringBuf[T, PT]) Capacity() int {
	return r.items.Capacity()
}

// AttachStorage to buffer
func (r *ringBuf[T, PT]) AttachStorage(s Storage[T, PT]) {
	r.items = s
}

// Put new data point to the buffer
func (r *ringBuf[T, PT]) Put(v *T) error {
	// 1. Get epoch since 1970 UTC
	pv := PT(v)
	epoch := pv.Timestamp().UTC().Unix()
	if epoch < 0 {
		return ErrOutOfRange
	}

	seq, idx := r.pos(epoch)
	pv.SetSequence(seq)
	if r.maxEpoch >= epoch {
		if item := r.items.At(idx); item != nil {
			// new data and existing data has equal timestamp,
			// merge data to existing record
			pv = PT(item)
			return pv.Merge(v)
		}
	} else {
		// We got new data, update maximum epoch and index
		r.maxEpoch = epoch
		r.maxIdx = idx
		r.latest = v
	}
	r.items.Put(idx, v)

	return nil
}

// Vacuum removes old items
func (r *ringBuf[T, PT]) Vacuum() error {
	if v := r.latest; v != nil {
		pv := PT(v)
		r.items.RemovesOlder(pv.Sequence() - 1)
	}
	return nil
}

// Each process items
func (r *ringBuf[T, PT]) Each(fn func(*T)) {
	if r.maxIdx != vInvalid {
		r.items.Each(r.maxIdx, fn)
	}
}

// Each process items
func (r *ringBuf[T, PT]) Walk(fn func(*T) bool) {
	if r.maxIdx != vInvalid {
		r.items.Walk(r.maxIdx, fn)
	}
}
