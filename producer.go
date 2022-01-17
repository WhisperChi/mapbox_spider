package main

import (
	"math"
)

type URLSchema struct {
	Type string
	Z    int
	X    int
	Y    int
}

type LimitRule struct {
	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64
	MinZ   int
	MaxZ   int
}

type URLItem struct {
	Z int
	X int
	Y int
}

func GenerateURLByRules(rule LimitRule, res chan<- URLItem) {
	for i := rule.MinZ; i <= rule.MaxZ; i++ {
		curSplitNumPerAxios := math.Pow(2, float64(i))
		// fmt.Println("split num is ", curSplitNumPerAxios)

		curDeltaX := 360 / curSplitNumPerAxios
		curDeltaY := 180 / curSplitNumPerAxios

		// fmt.Println("curDelta X,Y is zoom: ", i, " ", curDeltaX, ",", curDeltaY)

		curMinIndexX := int((rule.MinLon + 180.0) / curDeltaX)
		curMaxIndexX := int((rule.MaxLon + 180.0) / curDeltaX)

		curMinIndexY := int((rule.MinLat + 90.0) / curDeltaY)
		curMaxIndexY := int((rule.MaxLat + 90.0) / curDeltaY)

		// fmt.Println("min max lon index is zoom: ", i, "", curMinIndexX, " ", curMaxIndexX)
		// fmt.Println("min max lat index is zoom: ", i, " ", curMaxIndexX, " ", curMaxIndexX)

		var item URLItem
		item.Z = i
		item.X = curMinIndexX
		item.Y = curMinIndexY
		res <- item

		for j := curMinIndexX + 1; j < curMaxIndexX; j++ {
			for k := curMinIndexY + 1; k < curMaxIndexY; k++ {
				var item URLItem
				item.Z = i
				item.X = j
				item.Y = k
				res <- item
			}
		}
	}
}
