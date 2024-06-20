package main

import (
	"fmt"
	"math"
	"time"
)

func getYMDList(start int64, end int64) (ymdList []string) {
	if start > end {
		return
	}
	start = int64(math.Floor(float64(start)/(3600*24))) * 3600 * 24
	end = int64(math.Ceil(float64(end+1)/(3600*24))) * 3600 * 24
	for start < end {
		ymd := time.Unix(start, 0).Format("2006-01-02")
		ymdList = append(ymdList, ymd)
		start += 3600 * 24
	}
	return ymdList
}

func getTime(timeString string) (ts int64) {
	t, _ := time.Parse("2006-01-02T15:04:05+08:00", timeString)
	return t.Unix()
}

func main() {

	// 设置时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	time.Unix(1689183155, 0).In(location).Format("2006-01-02")
	fmt.Printf("1689183155=%v\n", time.Unix(1689183155, 0).In(location).Format("2006-01-02"))

	ttS, _ := time.Parse("2006-01-02 15:04:05", "2023-07-11 07:59:07")
	fmt.Printf("==== %v\n", ttS.Unix()-8*3600)
	// 1689033547

	var startList = []string{
		"2011-11-10T23:59:59+08:00",
		"2011-11-11T00:00:00+08:00",
		"2011-11-12T00:00:01+08:00",
		"2011-11-12T23:59:59+08:00",
		"2011-11-13T00:00:00+08:00",
		"2011-11-14T00:00:01+08:00",
	}
	for _, s := range startList {
		for _, e := range startList {
			st := getTime(s)
			et := getTime(e)
			if st > et {
				continue
			}
			fmt.Printf("%v   %v\n", s, e)
			fmt.Println(getYMDList(st, et))
		}
	}

	//go func() {
	//	ticker := time.NewTicker(3 * time.Second)
	//	for {
	//		select {
	//		case t := <-ticker.C:
	//			fmt.Println("Current time: ", t)
	//			time.Sleep(5 * time.Second)
	//		}
	//	}
	//}()
	//time.Sleep(10000 * time.Second)
}
