package manager

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestUser(t *testing.T) {
	TestingT(t)
}

type TestsSuite struct{}

var _ = Suite(&TestsSuite{})

func (s TestsSuite) Test_New(c *C) {
	//c.Skip("Not now")
	man := New("", "")
	c.Assert(man, NotNil)
}

func (s TestsSuite) Test_checkUrlsCount(c *C) {
	// c.Skip("Not now")
	man := New("", "")

	man.SetCountDownload(3)

	list := []string{"qwe1", "qwe2", "qwe3", "qwe4", "qwe", "qwe", "qwe", "qwe", "qwe", "qwe", "qwe", "qwe", "qwe", "qwe", "qwe", "qwe"}
	list2 := man.checkUrlsCount(list)

	c.Assert(list2, DeepEquals, []string{"qwe1", "qwe2", "qwe3"})
}
