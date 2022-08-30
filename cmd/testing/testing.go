package main

import (
	"fmt"
	"time"
)

var (
	tmZero = time.Time{}.UTC()
)

func tryTimeSeries(from, to time.Time, period, res time.Duration) {
	fmt.Println(from, to, from.Sub(tmZero).Seconds())
	from = from.UTC()
	to = to.UTC()

	const layout = time.RFC3339 //"2006-01-02 15:04:05Z"
	uFrom := from.Unix()
	uTo := to.Unix()

	sPeriod := int64(period.Seconds())
	sRes := int64(res.Seconds())
	szRes := sPeriod / sRes
	if sPeriod%sRes != 0 {
		szRes++
		sPeriod = szRes * sRes
	}

	fmt.Println(szRes, from, to, from.Sub(tmZero).Seconds())
	for tick := uFrom; tick < uTo; tick++ {
		tm := time.Unix(tick, 0).UTC()
		tmTrunc := tm.Truncate(res)
		tmStr := tm.Format(layout)
		tmTruncStr := tmTrunc.Format(layout)

		segment := tick / sPeriod
		idx := (tick % sPeriod) / sRes
		fmt.Printf("%s -> %s : %10d -> %d | %d\n",
			tmStr, tmTruncStr, segment, tick/szRes, idx)
	}

}
