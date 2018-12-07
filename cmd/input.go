package cmd

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"os"
)

// GetInputContent returns the content of input
func GetInputContent(input string) (io.Reader, error) {
	var r io.Reader

	u, err := url.Parse(input)
	if u.Scheme != "" {
		r, err = remoteList(input)
		if err != nil {
			return nil, err
		}
	} else if err != nil || u.Scheme == "" {
		r, err = localList(input)
		if err != nil {
			return nil, err
		}
	} else {
		r, err = defaultList()
		if err != nil {
			return nil, err
		}
	}

	var dst io.ReadWriter
	src := base64.NewDecoder(base64.StdEncoding, r)
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func localList(pathname string) (io.Reader, error) {
	f, err := os.Open(pathname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f, nil
}

func remoteList(urlpath string) (io.Reader, error) {
	resp, err := http.Get(urlpath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Body, nil
}

const gfwlist = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"

func defaultList() (io.Reader, error) {
	return remoteList(gfwlist)
}
