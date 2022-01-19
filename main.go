package main

import (
	"flag"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	CPUNum       *int
	MaxConsumers *int
	SaveDir      *string
	Type         *string
	Token        *string
	SKU          *string
	areaStr      *string
	MinLon       float64
	MaxLon       float64
	MinLat       float64
	MaxLat       float64
	MinZ         *int
	MaxZ         *int
)

func InitPara() {
	CPUNum = flag.Int("c", 1, "CPU num")
	MaxConsumers = flag.Int("maxc", 5, "Num of max consumers,this can speed up.")
	SaveDir = flag.String("d", "./default_download_path", "Path where you wan't to save data")
	Type = flag.String("t", "satellite", "satellite/street/terrain")
	Token = flag.String("token", "", "Your mapbox token.")
	SKU = flag.String("sku", "", "Your mapbox sku.")
	areaStr = flag.String("area", "-179.0,179.0,-89.0,89.0", "minLon,maxLon,minLat,maxLat")
	MinZ = flag.Int("minz", 0, "min zoom")
	MaxZ = flag.Int("maxz", 2, "max zoom")
}

func main() {

	// get para from cmd.
	InitPara()
	flag.Parse()

	// get para from config file
	var config Params
	GetParamsFromFile("./config.json", &config)

	if config.Token == "" && *Token == "" {
		fmt.Println("Please set your token.")
		return
	}
	if config.SKU == "" && *SKU == "" {
		fmt.Println(("Please set your sku."))
		return
	}

	var mpInfo MapboxInfo
	if *Token != "" {
		mpInfo.Token = *Token
	} else {
		mpInfo.Token = config.Token
	}

	if *SKU != "" {
		mpInfo.SKU = *SKU
	} else {
		mpInfo.SKU = config.SKU
	}

	mpInfo.SaveDir = *SaveDir

	aTmp := reflect.ValueOf(&config).Elem()
	mType := strings.Title(strings.ToLower(*Type))
	mDetails := aTmp.FieldByName(mType)

	mpInfo.Prefix = mDetails.FieldByName(strings.Title(strings.ToLower("prefix"))).String()
	mpInfo.Format = mDetails.FieldByName(strings.Title(strings.ToLower("format"))).String()
	mpInfo.FileSuffix = mDetails.FieldByName(strings.Title(strings.ToLower("filesuffix"))).String()

	runtime.GOMAXPROCS(*CPUNum)
	var wg sync.WaitGroup
	data := make(chan URLItem, 1000)

	for i := 0; i < *MaxConsumers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			DownloadItem(mpInfo, data)
		}()
	}

	strTmp := *areaStr
	strArr := strings.Split(strTmp, ",")
	var rule LimitRule

	rule.MinZ = *MinZ
	rule.MaxZ = *MaxZ
	if len(strArr) != 4 {
		// -179.0,179.0,-89.0,89.0
		rule.MinLon = -179.0
		rule.MaxLon = 179.0
		rule.MinLat = -89.0
		rule.MaxLat = 89.0
	} else {
		rule.MinLon, _ = strconv.ParseFloat(strArr[0], len(strArr[0]))
		rule.MaxLon, _ = strconv.ParseFloat(strArr[1], len(strArr[1]))
		rule.MinLat, _ = strconv.ParseFloat(strArr[2], len(strArr[2]))
		rule.MaxLat, _ = strconv.ParseFloat(strArr[3], len(strArr[3]))
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(data)

		GenerateURLByRules(rule, data)
	}()

	wg.Wait()
}
