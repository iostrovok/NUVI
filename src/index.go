package main

import (
	"flag"
	"fmt"
	"os"

	"extract_urls"
	"manager"
)

const (
	COUNT_GOROUTINE    = 8
	COUNT_DOWNLOAD     = -1
	URL                = "http://bitly.com/nuvi-plz"
	LIST_NAME          = "NEWS_XML"
	REDIS_HOST         = "localhost"
	REDIS_PORT         = 6379
	DEBUG              = false
	DOWNLOAD_OLD_FILES = false
)

var host, url, list_name string
var port, goes, files int
var debug, old bool

func main() {

	redis_url := fmt.Sprintf("%s:%d", host, port)

	// just for information
	fmt.Printf("Start redis. list: '%s',url: '%s'\n", list_name, redis_url)
	fmt.Printf("Remote url: %s\n", url)
	fmt.Printf("Count goroutine for dowlload: %d\n", goes)
	fmt.Printf("Count files for dowlload: %d\n", files)
	fmt.Printf("View debug info: %t\n", debug)
	fmt.Printf("Download old files: %t\n", old)

	// get list zip files
	file_urls, err := extract_urls.GetFileUrls(url)
	if len(file_urls) == 0 {
		print_result(fmt.Sprintf("No files for '%s'", URL), 0, 0)
		os.Exit(0)
	}

	if debug {
		fmt.Printf("Count files: '%d'\n", len(file_urls))
	}

	// init main nternal object
	man := manager.New(redis_url, list_name)

	man.DownloadOldFiles(old)   // Do you want to reload old files?
	man.SetDebug(debug)         // for debug
	man.SetCountDownload(files) // for limit downloading
	man.ManyLoaders(goes)

	// Download and proccess each file
	err = man.Make(file_urls)
	view := ""
	if err != nil {
		view = fmt.Sprintf("%s", err)
	}

	print_result(view, man.NewIds(), man.OldIds())

}

// print statistic
func print_result(err string, newIDs, oldIDs int) {
	if err != "" {
		fmt.Printf("ERROR: %s\n", err)
	}
	fmt.Printf("Downloaded. new: %d, old: %d\n", newIDs, oldIDs)
}

func init() {

	flag.BoolVar(&debug, "d", DEBUG, fmt.Sprintf("View debug info. Default is '%t'", DEBUG))
	flag.BoolVar(&old, "old", DOWNLOAD_OLD_FILES, fmt.Sprintf("Do you want to reload old files?. Default is '%t'", DOWNLOAD_OLD_FILES))

	flag.StringVar(&url, "url", URL, fmt.Sprintf("ZIP files page's url. Default is '%s'", URL))

	flag.StringVar(&list_name, "list", LIST_NAME, fmt.Sprintf("List name in Redis. Default is '%s'", LIST_NAME))
	flag.StringVar(&host, "host", REDIS_HOST, fmt.Sprintf("Redis host. Default is '%s'", REDIS_HOST))
	flag.IntVar(&port, "port", REDIS_PORT, fmt.Sprintf("Redis port. Default is %d", REDIS_PORT))

	flag.IntVar(&goes, "goes", COUNT_GOROUTINE, fmt.Sprintf("Count goroutine for dowlload. Default is %d", COUNT_GOROUTINE))
	flag.IntVar(&files, "files", COUNT_DOWNLOAD, fmt.Sprintf("Count zip for dowlload. Default is %d (no limit)", COUNT_DOWNLOAD))

	flag.Parse()
}
