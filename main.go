package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func whetherFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func writeDataToDisk(path string, data *[]byte) {
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

type zxyBaseInfo struct {
	BaseURL string
	EndURL  string
	Dir     string
	Suffix  string
}

type zxyProduct struct {
	Z int
	X int
	Y int
}

type downloadConfig struct {
	BaseURL string
	EndURL  string
	Dir     string
	Suffix  string
	MinZoom int
	MaxZoom int
}

type fileProduct struct {
	RelativePath string
	Data         []byte
}

func zxyProducer(id int, wg *sync.WaitGroup, minL int, maxL int, zxyChan chan zxyProduct) {

	defer wg.Done()
	defer close(zxyChan)

	for i := minL; i <= maxL; i++ {
		for j := 0; j < int(math.Pow(2, float64(i))); j++ {
			for k := 0; k < int(math.Pow(2, float64(i))); k++ {
				fmt.Printf("from zxy producer ## %d ## ,produce zxy data , z x y is %d %d %d\n", id, i, j, k)
				tmp := zxyProduct{Z: i, X: j, Y: k}
				zxyChan <- tmp
			}
		}
	}

}

func zxyConsumer(id int, wg *sync.WaitGroup, baseinfo zxyBaseInfo, jobs <-chan zxyProduct, results chan<- fileProduct) {
	defer wg.Done()
	defer close(results)

	fmt.Println("zxy consumer started")
	c := colly.NewCollector(colly.AllowURLRevisit())

	// c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 100})
	baseURL := baseinfo.BaseURL
	endURL := baseinfo.EndURL
	dataDir := baseinfo.Dir
	suffix := baseinfo.Suffix

	c.OnResponse(func(r *colly.Response) {

		data := r.Body
		relativePath := strings.Split(r.Request.URL.String(), baseURL)[1]
		relativePath = strings.Split(relativePath, endURL)[0]

		tmp := fileProduct{RelativePath: relativePath, Data: data}
		results <- tmp
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println(err)
	})

	for n := range jobs {
		z := n.Z
		x := n.X
		y := n.Y
		middleURL := strconv.Itoa(z) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y)
		url := baseURL + middleURL + endURL

		fmt.Printf("    from zxy consumer ## %d ##,handle zxy data %d,%d,%d\n", id, z, x, y)

		// if file not exist
		filePath := dataDir + "/" + middleURL + suffix
		if res := whetherFileExist(filePath); !res {
			c.Visit(url)
		} else {
			fmt.Println("file exist skip")
		}
	}
	fmt.Println("zxy consumer finished")

}

func fileConsumer(id int, wg *sync.WaitGroup, dir string, suffix string, jobs <-chan fileProduct) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("        from file worker #### %d ####, path is %s\n", id, j.RelativePath)
		path := dir + "/" + j.RelativePath + suffix
		writeDataToDisk(path, &j.Data)
	}
}

func startDownload(config downloadConfig, wg *sync.WaitGroup, zxyChan chan zxyProduct, fChan chan fileProduct) {

	defer wg.Done()

	baseURL := config.BaseURL
	endURL := config.EndURL
	suffix := config.Suffix
	dataDir := config.Dir
	minZ := config.MinZoom
	maxZ := config.MaxZoom

	baseInfo := zxyBaseInfo{BaseURL: baseURL, EndURL: endURL, Suffix: suffix, Dir: dataDir}

	// produce zxy data
	wg.Add(1)
	go zxyProducer(0, wg, minZ, maxZ, zxyChan)

	// file worker
	for w := 0; w < 5; w++ {
		wg.Add(1)
		go fileConsumer(w, wg, dataDir, suffix, fChan)
	}

	wg.Add(1)
	go zxyConsumer(0, wg, baseInfo, zxyChan, fChan)
}

func demDownloader(wg *sync.WaitGroup) {
	defer wg.Done()

	baseURL := "https://api.mapbox.com/raster/v1/mapbox.mapbox-terrain-dem-v1/"
	endURL := ".webp?sku=101wgnatEmnNk&access_token=<your-token>"
	suffix := ".webp"
	dataDir := "./dem"
	minZ := 0
	maxZ := 10
	zxyLength := 10
	var zxyChan = make(chan zxyProduct, zxyLength)
	var fChan = make(chan fileProduct)

	tmp := downloadConfig{BaseURL: baseURL, EndURL: endURL, Dir: dataDir, Suffix: suffix, MinZoom: minZ, MaxZoom: maxZ}
	wg.Add(1)
	go startDownload(tmp, wg, zxyChan, fChan)

	fmt.Println("test func end")
}

func main() {

	runtime.GOMAXPROCS(4)
	var wg sync.WaitGroup

	// dem downloader
	{
		wg.Add(1)
		go demDownloader(&wg)
	}

	wg.Wait()
}
