package main

import (
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(4)

	var wg sync.WaitGroup

	var dc DownloadConfig
	dc.SaveDir = "./data"

	var para Params
	GetParamsFromFile("./config.json", &para)

	var mpInfo MapboxInfo
	mpInfo.Token = para.Token
	mpInfo.Prefix = para.Satellite.Prefix
	mpInfo.SKU = para.SKU
	mpInfo.Format = para.Satellite.Format

	data := make(chan URLItem, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		DownloadItem(dc, mpInfo, data)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(data)

		var rule LimitRule
		rule.MinLon = -100
		rule.MaxLon = 100
		rule.MinLat = -70
		rule.MaxLat = 60
		rule.MinZ = 0
		rule.MaxZ = 2
		GenerateURLByRules(rule, data)
	}()

	wg.Wait()
}
