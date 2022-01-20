package main

import (
	"testing"
)

func TestGetParamsFromFile(t *testing.T) {
	var params Params
	var FILE_URL = "./config-example.json"
	GetParamsFromFile(FILE_URL, &params)
	if params.Token != "<tokenToBeReplaced>" {
		t.Error("getToken failed")
	}
}
