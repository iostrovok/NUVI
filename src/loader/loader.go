package loader

import (
	"sync"
	"sync/atomic"

	"redis"
)

type Loader struct {
	conn *redis.Client

	newID, oldID *int64 // for

	downloadOldFiles bool
}

// Constructor
func Start(conn *redis.Client, group *sync.WaitGroup, ch chan string, getdOldFiles bool, newIds, oldIds *int64) {
	this := &Loader{
		conn:  conn,
		newID: newIds,
		oldID: oldIds,

		downloadOldFiles: getdOldFiles,
	}

	group.Add(1)
	defer group.Done()

	for {
		select {
		case t, ok := <-ch:
			if !ok {
				return
			}
			err := this.Process(t)
			if err != nil {
				// TODO maka back chan and return error and return
				panic(err)
			}
		}
	}
}

func (this *Loader) Process(fileUrl string) error {

	if !this.downloadOldFiles {
		old, err := this.conn.CheckFile(fileUrl)
		if err != nil {
			return err
		}
		if old {
			return nil
		}
	}

	list, err := this.LoadAndUnzip(fileUrl)
	if err != nil {
		return err
	}

	err = this.putToRedis(list)
	if err != nil {
		return err
	}

	return this.conn.AddFile(fileUrl)
}

func (this *Loader) putToRedis(list []string) error {
	for _, news := range list {
		isNew, err := this.conn.Put(news)
		if err != nil {
			return err
		}

		if isNew {
			atomic.AddInt64(this.newID, 1)
		} else {
			atomic.AddInt64(this.oldID, 1)
		}
	}

	return nil
}
