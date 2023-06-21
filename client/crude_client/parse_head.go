package crude_client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func getEtag(rep *http.Response) (string, error) {
	etag := rep.Header.Get("etag")
	if etag == "" {
		return "", fmt.Errorf(
			"'etag' missing from response headers",
		)
	}
	return etag, nil
}

func getExpires(rep *http.Response) (time.Time, error) {
	datestring := rep.Header.Get("expires")
	if datestring == "" {
		return time.Time{}, fmt.Errorf(
			"'expires' missing from response headers",
		)
	}
	date, err := time.Parse(time.RFC1123, datestring)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"error parsing 'expires' header: %w",
			err,
		)
	}
	return date, nil
}

func getPages(rep *http.Response) (int, error) {
	pagesstring := rep.Header.Get("x-pages")
	if pagesstring == "" {
		return 0, fmt.Errorf(
			"'x-pages' missing from response headers",
		)
	}
	pages, err := strconv.Atoi(pagesstring)
	if err != nil {
		return 0, fmt.Errorf(
			"error parsing 'x-pages' header: %w",
			err,
		)
	}
	return pages, nil
}
