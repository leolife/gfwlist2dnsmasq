package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetInputContent returns the content of input
func GetInputContent(input string) (io.Reader, error) {
	var rc io.ReadCloser

	u, err := url.Parse(input)
	if u.Scheme != "" {
		rc, err = remoteList(input)
		if err != nil {
			return nil, err
		}
	} else if input != "" && u.Scheme == "" {
		rc, err = localList(input)
		if err != nil {
			return nil, err
		}
	} else {
		rc, err = defaultList()
		if err != nil {
			return nil, err
		}
	}

	dst := bytes.NewBufferString("")
	src := base64.NewDecoder(base64.StdEncoding, rc)
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return dst, nil
}

func localList(pathname string) (io.ReadCloser, error) {

	pathname = resolveHomeDir(pathname)

	pathname, err := filepath.Abs(pathname)
	f, err := os.Open(pathname)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func remoteList(urlpath string) (io.ReadCloser, error) {
	resp, err := http.Get(urlpath)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

const gfwlist = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"

func defaultList() (io.ReadCloser, error) {
	fmt.Println("Get default list")
	return remoteList(gfwlist)
}

func resolveHomeDir(pathname string) string {
	tempStr := ""

	usr, _ := user.Current()
	dir := usr.HomeDir

	if pathname == "~" {
		tempStr = dir
	} else if strings.HasPrefix(pathname, "~/") {
		tempStr = filepath.Join(dir, pathname[2:])
	}

	return tempStr
}
