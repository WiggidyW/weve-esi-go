package crude_client

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/WiggidyW/weve-esi/env"
)

const JSON_CONTENT_TYPE = "application/json"
const WWW_CONTENT_TYPE = "application/x-www-form-urlencoded"
const LOGIN_HOST = "login.eveonline.com"

func withHeadUserAgent(req *http.Request) {
	req.Header.Add("X-User-Agent", env.USER_AGENT)
}

func withHeadJsonContentType(req *http.Request) {
	req.Header.Add("Content-Type", JSON_CONTENT_TYPE)
}

func withHeadWwwContentType(req *http.Request) {
	req.Header.Add("Content-Type", WWW_CONTENT_TYPE)
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

func withHeadLoginHost(req *http.Request) {
	req.Header.Add("Host", LOGIN_HOST)
}
