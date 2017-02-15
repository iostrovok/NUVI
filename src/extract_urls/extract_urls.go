package extract_urls

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

var getHrefReg = regexp.MustCompile(`<a href="([^"]+.zip)">`)

// <a href="1486756614459.zip">
// http://bitly.com/nuvi-plz

func GetFileUrls(url string) ([]string, error) {
	finalURL, body, err := getUrls(url)
	if err != nil {
		return nil, err
	}

	files := parse(body)
	for i := range files {
		files[i] = finalURL + files[i]
	}

	return files, nil
}

func getUrls(url string) (string, string, error) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return nil },
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", "", err
	}

	finalURL := resp.Request.URL.String()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	return finalURL, string(body), nil
}

func clean_url(url string) string {
	if url == "" {
		return ""
	}

	url = strings.TrimRight(url, "/")
	return url + "/"
}

func parse(html string) []string {
	list := getHrefReg.FindAllStringSubmatch(html, -1)

	out := []string{}
	for _, l := range list {
		if len(l) == 2 {
			out = append(out, l[1])
		}
	}

	sort.Strings(out)

	return out
}
