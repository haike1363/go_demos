package main

import (
	"fmt"
	"math"
)

func floorByRound(start float64, round int64, up bool) float64 {
	result := int64(1)
	for ; start > float64(round); {
		start /= float64(round)
		result *= round
	}
	if up {
		return float64(result) * math.Ceil(start)
	} else {
		return float64(result) * math.Floor(start)
	}
}

func getStorageUnit(val float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	i := 0
	for ; i < len(units) && val >= 1024; i++ {
		val /= 1024
	}
	if i == len(units) {
		i--
	}
	return fmt.Sprintf("%v%v", int64(val), units[i])
}

type MinMaxValue struct {
	Min float64
	Max float64
}

func main() {
	minMaxValue := &MinMaxValue{
		Min: 0,
		Max: 3.16 * 1000 * 1000 * 1000,
	}

	genFilterFunc := func(round int64, getUnitFunc func(val float64) string) {
		minMaxValue.Min = math.Floor(minMaxValue.Min)
		minMaxValue.Max = math.Ceil(minMaxValue.Max)
		deltaValue := (minMaxValue.Max - minMaxValue.Min) / float64(8)
		if deltaValue < float64(round) {
			deltaValue = float64(round)
		}
		for start := floorByRound(minMaxValue.Min, round, false); start < minMaxValue.Max; {
			end := floorByRound(start+deltaValue, round, true)
			fmt.Printf("%v~%v\n", getUnitFunc(start), getUnitFunc(end))
			if start >= end {
				break
			}
			start = end
		}
	}

	genFilterFunc(1024, getStorageUnit)
}
