package timeseries_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ipsusila/timeseries"
	"github.com/stretchr/testify/assert"
)

type Geodata struct {
	timeseries.Datum
	val int64
}

func (g *Geodata) Merge(o *Geodata) error {
	return nil
}

func TestRingBuffer(t *testing.T) {
	buf, err := timeseries.NewBuffer[Geodata](2*time.Hour, 2*time.Second, timeseries.MemStore)
	assert.NoError(t, err)
	lt := buf.Latest()
	assert.Nil(t, lt)

	gd := Geodata{
		Datum: timeseries.Datum{
			Ts: time.Now(),
		},
		val: 10,
	}
	buf.Put(&gd)
	lt = buf.Latest()
	assert.NotNil(t, lt)
}

func TestFillBuffer(t *testing.T) {
	buf, err := timeseries.NewBuffer[Geodata](2*time.Hour, 2*time.Second, timeseries.MemStore)
	assert.NoError(t, err)
	now := time.Now()
	base := now.Add(-90 * time.Minute)
	src := rand.NewSource(now.UnixNano())
	r := rand.New(src)
	n := 10000
	lim := int64(180 * time.Minute)
	for i := 0; i < n; i++ {
		rv := time.Duration(r.Int63n(lim))
		ts := base.Add(rv)
		gd := Geodata{
			Datum: timeseries.Datum{
				Ts: ts,
			},
			val: int64(i),
		}
		buf.Put(&gd)
	}

	fd, err := os.Create("testout.txt")
	assert.NoError(t, err)
	defer fd.Close()
	var first *Geodata
	fo := func(v *Geodata) {
		delta := v.Ts.Sub(first.Ts)
		fmt.Fprintf(fd, "%s | %05d | %6.0f | %d\n", v.Ts.Format(time.RFC3339), v.val, delta.Seconds(), v.Seq)
	}
	buf.Each(func(v *Geodata) {
		if first == nil {
			first = v
		}
		fo(v)
	})
	fmt.Fprintln(fd, "=====")
	fo(buf.Latest())
	fmt.Fprintln(fd, "=====")

	buf.Walk(func(v *Geodata) bool {
		if v.Ts.After(now) {
			return false
		}
		fo(v)
		return true
	})
}

func getBuffer[T any, PT timeseries.Timestamper[T]](period, res time.Duration) timeseries.Buffer[T, PT] {
	buf, err := timeseries.NewBuffer[T, PT](period, res, timeseries.MemStore)
	if err != nil {
		return nil
	}
	return buf
}

func BenchmarkRingBuffer1Day(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getBuffer[Geodata](24*time.Hour, time.Second)
	}
}

func BenchmarkRingBuffer1Week(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getBuffer[Geodata](7*24*time.Hour, time.Second)
	}
}
