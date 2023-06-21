package crude_client

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/WiggidyW/weve-esi/env"
)

const CONTENT_TYPE = "application/json"

func withHeadUserAgent(req *http.Request) {
	req.Header.Add("X-User-Agent", env.USER_AGENT)
}

func withHeadContentType(req *http.Request) {
	req.Header.Add("Content-Type", CONTENT_TYPE)
}

func withHeadAuthorization(req *http.Request, auth string) {
	if auth != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth))
	}
}

func withHeadEtag(req *http.Request, etag string) {
	if etag != "" {
		req.Header.Add("If-None-Match", etag)
	}
}

func withHeadAuthAuthorization(req *http.Request) {
	basic_auth := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", env.CLIENT_ID, env.CLIENT_SECRET),
	))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basic_auth))
}
