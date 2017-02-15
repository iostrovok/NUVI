package loader

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"net/http"
)

// dowload and unzip file
func (this *Loader) LoadAndUnzip(url string) ([]string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// ReadAll reads from readCloser until EOF and returns the data as a []byte
	b, err := ioutil.ReadAll(resp.Body) // The readCloser is the one from the zip-package
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// bytes.Reader implements io.Reader, io.ReaderAt, etc. All you need!
	readerAt := bytes.NewReader(b)

	z, err := zip.NewReader(readerAt, int64(len(b)))
	if err != nil {
		return nil, err
	}

	out := []string{}

	// Iterate through the files in the archive,
	// process contents.
	for _, f := range z.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(rc)
		out = append(out, buf.String())
		rc.Close()
	}

	return out, nil
}
