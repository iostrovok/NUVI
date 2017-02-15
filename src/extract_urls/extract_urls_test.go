package extract_urls

import (
	. "gopkg.in/check.v1"
	"testing"
)

/*
 */

var links []string = []string{"1486756614459.zip", "1486756815765.zip", "1486756956243.zip", "1486757319200.zip", "1486757667279.zip", "1486758003597.zip",
	"1486758485600.zip", "1486758834945.zip",
}

var test_html string = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
<html>
 <head>
  <title>Index of /5Rh5AMTrc4Pv/mainstream/posts</title>
 </head>
 <body>
<h1>Index of /5Rh5AMTrc4Pv/mainstream/posts</h1>
<table><tr><th><a href="?C=N;O=D">Name</a></th><th><a href="?C=M;O=A">Last modified</a></th><th><a href="?C=S;O=A">Size</a></th><th><a href="?C=D;O=A">Description</a></th></tr><tr><th colspan="4"><hr></th></tr>
<tr><td><a href="/5Rh5AMTrc4Pv/mainstream/">Parent Directory</a></td><td>&nbsp;</td><td align="right">  - </td><td>&nbsp;</td></tr>
<tr><td><a href="1486756614459.zip">1486756614459.zip</a></td><td align="right">10-Feb-2017 22:00  </td><td align="right">9.9M</td><td>&nbsp;</td></tr>
<tr><td><a href="1486756815765.zip">1486756815765.zip</a></td><td align="right">10-Feb-2017 22:02  </td><td align="right">9.9M</td><td>&nbsp;</td></tr>
<tr><td><a href="1486756956243.zip">1486756956243.zip</a></td><td align="right">10-Feb-2017 22:08  </td><td align="right"> 10M</td><td>&nbsp;</td></tr>
<tr><td><a href="1486757319200.zip">1486757319200.zip</a></td><td align="right">10-Feb-2017 22:14  </td><td align="right"> 10M</td><td>&nbsp;</td></tr>
<tr><td><a href="1486757667279.zip">1486757667279.zip</a></td><td align="right">10-Feb-2017 22:20  </td><td align="right">9.9M</td><td>&nbsp;</td></tr>
<tr><td><a href="1486758003597.zip">1486758003597.zip</a></td><td align="right">10-Feb-2017 22:28  </td><td align="right">9.9M</td><td>&nbsp;</td></tr>
<tr><td><a href="1486758485600.zip">1486758485600.zip</a></td><td align="right">10-Feb-2017 22:33  </td><td align="right">9.9M</td><td>&nbsp;</td></tr>
<tr><td><a href="1486758834945.zip">1486758834945.zip</a></td><td align="right">10-Feb-2017 22:37  </td><td align="right">9.9M</td><td>&nbsp;</td></tr>
<tr><td>
`

func TestUser(t *testing.T) {
	TestingT(t)
}

type TestsSuite struct{}

var _ = Suite(&TestsSuite{})

func (s TestsSuite) Test_parse_href(c *C) {
	//c.Skip("Not now")
	list := parse(test_html)
	c.Assert(list, NotNil)
	c.Assert(list, DeepEquals, links)
}

func (s TestsSuite) Test_clean_url(c *C) {
	//c.Skip("Not now")
	href := clean_url("")
	c.Assert(href, Equals, "")

	href = clean_url("/")
	c.Assert(href, Equals, "/")

	href = clean_url("sdasdsad/")
	c.Assert(href, Equals, "sdasdsad/")

	href = clean_url("sdasdsad")
	c.Assert(href, Equals, "sdasdsad/")
}
