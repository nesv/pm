package pm

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var FetchSupportedSchemes = []string{"http", "https", "file"}

type FetchFunc func(*url.URL) (io.ReadCloser, error)

func fetchHTTP(u *url.URL) (io.ReadCloser, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("error: failed to fetch %q: %v", u.String(), err)
	}

	return resp.Body, nil
}

func fetchLocalFile(u *url.URL) (io.ReadCloser, error) {
	src, err := os.Open(u.Path)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("error: cannot fetch %q: does not exist", u.String())
	} else if err != nil {
		return nil, fmt.Errorf("error: cannot fetch %q: %v", u.String(), err)
	}

	return src, nil
}

func Fetch(urlStr string) (io.ReadCloser, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var fetch FetchFunc

	switch u.Scheme {
	case "http", "https":
		fetch = fetchHTTP
	case "", "file":
		fetch = fetchLocalFile
	default:
		return nil, fmt.Errorf("pm: unsupported scheme in URL: %v", urlStr)
	}

	return fetch(u)
}
