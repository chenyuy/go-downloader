package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func downloadFile(fileURL string) {
	parsed, err := url.Parse(fileURL)
	if err != nil {
		log.Println(err)
		wg.Done()
		return
	}
	paths := strings.Split(parsed.Path, "/")
	filename := parsed.Hostname()
	if len(paths) > 0 {
		filename = paths[len(paths)-1]
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		wg.Done()
		return
	}
	defer f.Close()
	log.Printf("Download from %s to %s\n", fileURL, filename)
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Println(err)
		wg.Done()
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Println(err)
		wg.Done()
		return
	}
	wg.Done()
}

func downloadFiles(urls []string, batch int) {
	for i := 0; i < len(urls); i += batch {
		for j := 0; i+j < len(urls) && j < batch; j++ {
			wg.Add(1)
			go downloadFile(urls[i+j])
		}
		wg.Wait()
	}
}

func main() {
	filename := flag.String("f", "", "path to a file containing a list of urls to download (one per line)")
	batch := flag.Int("b", 20, "maximum number of concurrent downloads (default: 20)")
	flag.Parse()
	if filename == nil || *filename == "" || batch == nil {
		flag.Usage()
		return
	}
	f, err := os.Open(*filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	downloadFiles(strings.Split(string(content), "\n"), *batch)
}
