package main

import (
	"fmt"
	"time"
)

func getYearAndMonth(dd time.Time) (start time.Time, end time.Time) {
	year, month, _ := dd.Date()
	loc := dd.Location()

	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	return startOfMonth, endOfMonth
}

func GetStartOfNextMonth(dd time.Time) (start time.Time) {
	year, month, _ := dd.Date()
	loc := dd.Location()
	return time.Date(year, month+1, 1, 0, 0, 0, 0, loc)
}

func main() {
	now := time.Now().AddDate(0, 1, -1)
	fmt.Println(now)
	a, b := getYearAndMonth(now)
	fmt.Println(a, b)

	fmt.Println(GetStartOfNextMonth(now))
}
