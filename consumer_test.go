package main

import (
	"sync"
	"testing"
)

func TestDownloadItem(t *testing.T) {
	var mpInfo MapboxInfo
	mpInfo.Token = "xxxToken"
	mpInfo.Prefix = "https://api.mapbox.com/v4/mapbox.satellite/"
	mpInfo.Format = ".webp"
	mpInfo.SKU = "sssSDK"
	data := make(chan URLItem, 10)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		DownloadItem(mpInfo, data)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
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
