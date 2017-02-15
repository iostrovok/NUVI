package redis

import (
	"crypto/md5"
	"fmt"

	red "github.com/garyburd/redigo/redis"
)

const LIST_NAME = "NEWS_XML"
const HASH_NAME = "NEWS_XML_HASH"
const FILES_NAME = "NEWS_XML_FILES"
const PROTOCOL_SCHEME = "tcp"

type Client struct {
	// client *red.Client
	client            red.Conn
	list, hash, files string
}

// host_port likes "localhost:6379"
func New(host_port string, list ...string) (*Client, error) {

	c := &Client{
		list:  LIST_NAME,
		hash:  HASH_NAME,
		files: FILES_NAME,
	}

	// for test
	if len(list) > 1 {
		c.list = list[0]
		c.hash = list[0] + "_HASH"
		c.files = list[0] + "_FILES"
	}

	cl, err := red.Dial(PROTOCOL_SCHEME, host_port)
	if err != nil {
		return nil, err
	}

	c.client = cl

	return c, nil
}

// if we don't have item we put it to Redis
func (this *Client) Put(item string) (bool, error) {
	key := fmt.Sprintf("%x", md5.Sum([]byte(item)))
	return this.findOrInsertInList(item, key)
}

// check downloaded file
func (this *Client) CheckFile(val string) (bool, error) {
	val = fmt.Sprintf("%x", md5.Sum([]byte(val)))
	return this.HExists(this.files, val)
}

// set downloaded file
func (this *Client) AddFile(val string) error {
	val = fmt.Sprintf("%x", md5.Sum([]byte(val)))
	return this.HSet(this.files, val)
}

/*
*****
	Internal function
*****
*/
// Search and/or insert LUA script
var findOrInsertInListSetScript = red.NewScript(1, `
	  if redis.call( 'HEXISTS', ARGV[3], ARGV[4] ) == 1 then
             return 1
      end
      redis.call( 'LPUSH', ARGV[1], ARGV[2] )
      redis.call( 'HSET', ARGV[3], ARGV[4], 1 )
      return -1
`)

// Search LUA script
var searchItemInListScript = red.NewScript(1, `
	local n = redis.call('LLEN', ARGV[1]) - 1
      for i=0,n do
         if redis.call( 'LINDEX', ARGV[1], i ) == ARGV[2] then
             return i
         end
      end
      return -1
`)

// finds existed items in list by LUA script
func (this *Client) findOrInsertInList(item, key string) (bool, error) {

	v, err := findOrInsertInListSetScript.Do(this.client, 1, this.list, item, this.hash, key)

	if err != nil {
		return false, err
	}

	if i, ok := v.(int64); ok {
		return i < 0, nil
	}

	return false, nil
}

// Seach exists items in list with LUA script
func (this *Client) searchOld(item string) bool {
	v, err := searchItemInListScript.Do(this.client, 1, this.list, fmt.Sprintf("%s", item))

	if err != nil {
		return false
	}

	if i, ok := v.(int64); ok {
		return i > -1
	}

	return false
}

func (this *Client) Del() {
	this.client.Do("DEL", this.hash)
	this.client.Do("DEL", this.files)
	this.client.Do("DEL", this.list)
}

func (this *Client) HSet(key, val string) error {
	_, err := this.client.Do("HSET", key, val, 1)
	return err
}

func (this *Client) HExists(key, val string) (bool, error) {
	v, err := this.client.Do("HEXISTS", key, val)
	if err != nil {
		return false, err
	}

	if i, ok := v.(int64); ok {
		return i != 0, nil
	}

	return false, nil
}

func (this *Client) LPush(item string) error {
	return this.client.Send("LPUSH", this.list, item)
}

func (this *Client) LRange(from, to int) ([]string, error) {

	out := []string{}

	reply, err := red.Values(this.client.Do("LRANGE", this.list, from, to))
	if err != nil {
		return nil, err
	}

	for _, b := range reply {
		out = append(out, fmt.Sprintf("%s", b))
	}

	return out, nil
}
