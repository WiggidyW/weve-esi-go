package response

import (
	"time"
)

type EsiHeadResponse struct {
	Pages   int
	Expires time.Time
}

func (r *EsiHeadResponse) GetPages() int {
	return r.Pages
}

type EsiResponse struct {
	Json    []map[string]interface{}
	Etag    string
	Expires time.Time
}
