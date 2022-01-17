package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type PartURL struct {
	Prefix string
	Format string
}

type Params struct {
	Token  string
	Styles []string
	SKU    string
	// Prefix string
	// Format string
	Satellite PartURL
	Street    PartURL
	Terrain   PartURL
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

func GetParamsFromCmd(cmd string, params *CmdParams) error {
	split := strings.Split(cmd, "-")
	cmdParamsMap := make(map[string]string)

	if len(split)%2 != 0 {
		return errors.New("num of params error")
	}
	for i := 0; i < len(split); i++ {
		cmdParamsMap[split[i]] = split[1]
	}

	fmt.Println(split)

	return nil
}
