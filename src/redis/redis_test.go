package redis

import (
	. "gopkg.in/check.v1"
	"testing"
)

const LIST_NAME_TEST = "NEWS_XML_TEST"
const TEST_REDIS = "localhost:6379"

/*
 */

func Test(t *testing.T) {
	TestingT(t)
}

type TestsSuite struct{}

var _ = Suite(&TestsSuite{})

func (s TestsSuite) Test_New(c *C) {
	// c.Skip("Not now")
	red, err := New(TEST_REDIS, LIST_NAME_TEST)
	c.Assert(err, IsNil)
	c.Assert(red, NotNil)
}

/*
redis>  HSET myhash field1 "foo"
(integer) 1
redis>  HEXISTS myhash field1
(integer) 1
redis>  HEXISTS myhash field2
(integer) 0
*/

func (s TestsSuite) Test_HSet_HExists(c *C) {
	// c.Skip("Not now")
	red, err := New(TEST_REDIS, LIST_NAME_TEST)
	c.Assert(err, IsNil)
	red.Del()

	key := red.files

	red.HSet(key, "world")
	red.HSet(key, "hello")

	find, err := red.HExists(key, "world")
	c.Assert(err, IsNil)
	c.Assert(find, DeepEquals, true)

	find, err = red.HExists(key, "hello")
	c.Assert(err, IsNil)
	c.Assert(find, DeepEquals, true)

	find, err = red.HExists(key, "bad_test")
	c.Assert(err, IsNil)
	c.Assert(find, DeepEquals, false)

	red.Del()

}
func (s TestsSuite) Test_AddFile_CheckFile(c *C) {
	// c.Skip("Not now")
	red, err := New(TEST_REDIS, LIST_NAME_TEST+"_")
	c.Assert(err, IsNil)
	red.Del()

	red.AddFile("world")
	red.AddFile("hello")

	find, err := red.CheckFile("world")
	c.Assert(err, IsNil)
	c.Assert(find, DeepEquals, true)

	find, err = red.CheckFile("hello")
	c.Assert(err, IsNil)
	c.Assert(find, DeepEquals, true)

	find, err = red.CheckFile("bad_test")
	c.Assert(err, IsNil)
	c.Assert(find, DeepEquals, false)

	// red.Del()
}

/*
redis>  LPUSH mylist "world"
redis>  LPUSH mylist "hello"
redis>  LRANGE mylist 0 -1
1) "hello"
2) "world"
*/
func (s TestsSuite) Test_LPUSH_LRANGE(c *C) {
	// c.Skip("Not now")
	red, err := New(TEST_REDIS, LIST_NAME_TEST)
	c.Assert(err, IsNil)
	red.Del()

	red.LPush("world")
	red.LPush("hello")

	list, err := red.LRange(0, -1)
	red.Del()
	c.Assert(list, DeepEquals, []string{"hello", "world"})
	c.Assert(err, IsNil)
}

func (s TestsSuite) Test_searchOld(c *C) {
	// c.Skip("Not now")
	red, err := New(TEST_REDIS, LIST_NAME_TEST)
	c.Assert(err, IsNil)
	red.Del()

	find := red.searchOld("world")
	c.Assert(find, Equals, false)

	red.LPush("world")
	red.LPush("hello")

	find = red.searchOld("world")
	c.Assert(find, Equals, true)

	red.Del()
}

func (s TestsSuite) Test_Put(c *C) {
	// c.Skip("Not now")
	red, err := New(TEST_REDIS, LIST_NAME_TEST)
	c.Assert(err, IsNil)
	red.Del()

	isNew, err := red.Put("world")
	c.Assert(err, IsNil)
	c.Assert(isNew, Equals, true)

	isNew, err = red.Put("world")
	c.Assert(err, IsNil)
	c.Assert(isNew, Equals, false)

	find := red.searchOld("world")
	c.Assert(find, Equals, true)

	list, err := red.LRange(0, -1)
	c.Assert(list, DeepEquals, []string{"world"})

	red.Del()
}
