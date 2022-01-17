package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestGenerateURLByRules(t *testing.T) {
	var rule LimitRule
	rule.MinLon = -100
	rule.MaxLon = 100
	rule.MinLat = -70
	rule.MaxLat = 60
	rule.MinZ = 0
	rule.MaxZ = 2
	res := make(chan URLItem, 10)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(res)
		for n := range res {
			fmt.Println("in chan URLItem is ", n.Z, " ", n.X, " ", n.Y)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		GenerateURLByRules(rule, res)
	}()

	wg.Wait()

}
