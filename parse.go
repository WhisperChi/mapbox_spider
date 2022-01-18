package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type PartURL struct {
	Prefix     string `json:"prefix"`
	Format     string `json:"format"`
	FileSuffix string `json:"filesuffix"`
}

type Params struct {
	Token     string  `json:"token"`
	SKU       string  `json:"sku"`
	Satellite PartURL `json:"satellite"`
	Street    PartURL `json:"street"`
	Terrain   PartURL `json:"filesuffix"`
}

type CmdParams struct {
	CpuNum   uint
	MaxPipes uint
}

func GetParamsFromFile(file string, params *Params) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), params)

	defer jsonFile.Close()
}
