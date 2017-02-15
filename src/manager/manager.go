package manager

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"loader"
	"redis"
)

type Manager struct {
	group            *sync.WaitGroup
	new_task         chan string
	debug            bool
	downloadOldFiles bool

	redis_url, redis_list string

	countDownload int

	oldIds, newIds *int64 // for checking of count of old/new news in Redis
}

// Constructor
func New(redis_url, redis_list string) *Manager {
	var oldIds int64 = 0
	var newIds int64 = 0

	M := &Manager{
		group:         &sync.WaitGroup{},
		new_task:      make(chan string, 120),
		countDownload: -1,
		oldIds:        &oldIds,
		newIds:        &newIds,
		redis_url:     redis_url,
		redis_list:    redis_list,
	}

	return M
}

func (this *Manager) SetDebug(t bool) {
	this.debug = t
}

func (this *Manager) DownloadOldFiles(t bool) {
	this.downloadOldFiles = t
}

//
func (this *Manager) AddLoader() {

	// init new DB connection
	conn, err := redis.New(this.redis_url, this.redis_list)
	if err != nil {
		fmt.Printf("ERROR redis_url: %s, redis_list: %s\n", this.redis_url, this.redis_list)
		panic(err)
	}

	go loader.Start(conn, this.group, this.new_task, this.downloadOldFiles, this.newIds, this.oldIds)
}

// Start n subs
func (this *Manager) ManyLoaders(goroutines int) {
	for i := 0; i < goroutines; i++ {
		this.AddLoader()
	}
}

func (this *Manager) OldIds() int {
	return int(atomic.LoadInt64(this.oldIds))
}

func (this *Manager) NewIds() int {
	return int(atomic.LoadInt64(this.newIds))
}

func (this *Manager) SetCountDownload(count_files int) {
	this.countDownload = count_files
}

// Check count download. Need for tests.
func (this *Manager) checkUrlsCount(urls []string) []string {

	if this.countDownload < 0 || len(urls) < this.countDownload {
		return urls
	}

	return urls[:this.countDownload]
}

// Procces list of urls
func (this *Manager) Make(urls []string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%s", r))
		}
	}()

	urls = this.checkUrlsCount(urls)

	for _, url := range urls {
		if this.debug {
			fmt.Printf("url: %s\n", url)
		}
		this.new_task <- url
	}

	close(this.new_task)

	this.group.Wait()

	return
}
