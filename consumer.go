package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type DownloadConfig struct {
	SaveDir string
}

type MapboxInfo struct {
	Token  string
	SKU    string
	Prefix string
	Format string
}

func DownloadItem(config DownloadConfig, mapboxInfo MapboxInfo, data <-chan URLItem) {
	c := colly.NewCollector(colly.AllowURLRevisit())
	baseURL := mapboxInfo.Prefix
	extraParams := "?sku=" + mapboxInfo.SKU + "&access_token=" + mapboxInfo.Token

	var wg sync.WaitGroup

	c.OnResponse(func(r *colly.Response) {
		// data := r.Body
		fmt.Println(r.Request.URL.String())

		data := r.Body
		relativePath := strings.Split(r.Request.URL.String(), baseURL)[1]
		relativePath = strings.Split(relativePath, mapboxInfo.Format+extraParams)[0]
		path := config.SaveDir + "/" + relativePath + mapboxInfo.Format
		wg.Add(1)
		go func() {
			defer wg.Done()
			FileWriter(path, &data)
		}()
		wg.Wait()
	})

	c.OnError(func(r *colly.Response, err error) {
		// TODO: save state
		// relativePath := strings.Split(r.Request.URL.String(), baseURL)[1]
		// relativePath = strings.Split(relativePath, mapboxInfo.Format+extraParams)[0]
		fmt.Println("Error,  ", err, " path is ", r.Request.URL.String())
	})

	for n := range data {
		z := n.Z
		x := n.X
		y := n.Y

		para := strconv.Itoa(z) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + mapboxInfo.Format
		url := baseURL + para + extraParams
		c.Visit(url)
	}
}

func FileWriter(path string, data *[]byte) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("path %s exist\n", path)
		return
	}

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := f.Write(*data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
	f.Sync()
}
