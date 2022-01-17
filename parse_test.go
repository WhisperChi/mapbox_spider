package main

import (
	"testing"
)

func TestGetParamsFromFile(t *testing.T) {
	var params Params
	var FILE_URL = "../config-example.json"
	GetParamsFromFile(FILE_URL, &params)
	if params.Token != "<tokenToBeReplaced>" {
		t.Error("getToken failed")
	}
}

func TestGetParamsFromCmd(t *testing.T) {
	var cmdParams CmdParams
	var str = "-c 4 -j 8"
	GetParamsFromCmd(str, &cmdParams)
}

func TestGetParamsFromCmdWithErrorParams(t *testing.T) {
	var cmdParams CmdParams
	var str = "-c 4 -j "
	err := GetParamsFromCmd(str, &cmdParams)

	if err != nil {
		t.Log("Cmd's num error.")
	} else {
		t.Error("The error of cmd's num can't be found. ")
	}
}
