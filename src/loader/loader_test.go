package loader

import (
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"redis"
)

const LIST_NAME_TEST = "NEWS_XML_TEST"
const TEST_REDIS = "localhost:6379"

func TestLoader(t *testing.T) {
	TestingT(t)
}

type TestsSuite struct{}

var _ = Suite(&TestsSuite{})

func testLoaoder() *Loader {
	var t1 int64 = 0
	var t2 int64 = 0

	return &Loader{
		conn:  nil,
		newID: &t1,
		oldID: &t2,
	}

}

func (s TestsSuite) Test_LoadAndUnzip(c *C) {
	// c.Skip("Not now")
	ts := httptest.NewServer(http.FileServer(http.Dir("../../")))
	defer ts.Close()

	myloader := testLoaoder()

	// list, err := load.LoadAndUnzip("http://feed.omgili.com/5Rh5AMTrc4Pv/mainstream/posts/1487021871192.zip")
	list, err := myloader.LoadAndUnzip(ts.URL + "/1487021871192.zip")

	c.Assert(err, IsNil)
	c.Assert(len(list), Equals, 2417)
}

func (s TestsSuite) Test_PutToRedis(c *C) {
	// c.Skip("Not now")

	conn, err := redis.New(TEST_REDIS, LIST_NAME_TEST)
	c.Assert(err, IsNil)

	conn.Del()

	myloader := testLoaoder()
	myloader.conn = conn

	err = myloader.putToRedis([]string{"qqqq", "qqqq1", "qqqq2"})
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(*myloader.newID, Equals, int64(3))
	c.Assert(*myloader.oldID, Equals, int64(0))

	err = myloader.putToRedis([]string{"qqqq", "qqqq1", "qqqq2", "qqqq3"})
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(*myloader.newID, Equals, int64(4))
	c.Assert(*myloader.oldID, Equals, int64(3))

	// clean redis
	conn.Del()
}
