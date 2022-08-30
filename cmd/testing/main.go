package main

import "time"

func main() {
	from := time.Now().Truncate(time.Hour * 24)
	to := from.Add(time.Hour * 24)

	//tryTimeSeries(from, to, 24*time.Hour, 1*time.Second)
	tryTimeSeries(from, to, 17*time.Second, 3*time.Second)
}
