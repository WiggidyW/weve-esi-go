package crude_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/WiggidyW/weve-esi/client/crude_client/cache"
	"github.com/WiggidyW/weve-esi/client/crude_client/response"
	"github.com/WiggidyW/weve-esi/env"
)

const AUTH_URL = "https://login.eveonline.com/v2/oauth/token"

type CrudeClient struct {
	Cache  cache.Cache
	Client *http.Client
}

func NewCrudeClient(run_local bool) *CrudeClient {
	var c cache.Cache
	if run_local {
		c = cache.NewLocalCache()
	} else {
		c = cache.NewCache()
	}
	return &CrudeClient{
		Cache: c,
		Client: &http.Client{
			Timeout: env.CLIENT_TIMEOUT,
		},
	}
}

func (c *CrudeClient) RequestNoCache(
	ctx context.Context,
	url string,
	method string,
	auth string,
) (*response.EsiResponse, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, RequestParamsError{err}
	}

	withHeadUserAgent(req)
	withHeadJsonContentType(req)
	withHeadAuthorization(req, auth)

	srvr_rep, err := c.Client.Do(req)
	if srvr_rep != nil {
		defer srvr_rep.Body.Close()
	}
	if err != nil {
		return nil, HttpError{err}
	}

	status := srvr_rep.StatusCode
	if status != http.StatusOK {
		return nil, newStatusError(srvr_rep)
	}

	nocache_rep := &response.EsiResponse{}
	err = json.NewDecoder(srvr_rep.Body).Decode(&nocache_rep.Json)
	if err != nil {
		return nil, MalformedResponse{fmt.Errorf(
			"error decoding response body as json: %w",
			err,
		)}
	}

	return nocache_rep, nil
}

func (c *CrudeClient) Request(
	ctx context.Context,
	url string,
	method string,
	auth string,
) (*response.EsiResponse, error) {
	c.Cache.Lock(url)
	defer c.Cache.Unlock(url)

	cache_rep, err := c.Cache.Get(ctx, url)
	if err != nil {
		return nil, CacheGetError{err}
	}
	if cache_rep != nil {
		now := time.Now()
		if now.Before(cache_rep.Expires) {
			return cache_rep, nil
		}
	} else {
		cache_rep = &response.EsiResponse{}
	}

	fmt.Printf("Sending '%s' request to '%s'\n", method, url)

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, RequestParamsError{err}
	}

	withHeadUserAgent(req)
	withHeadJsonContentType(req)
	withHeadAuthorization(req, auth)
	withHeadEtag(req, cache_rep.Etag)

	srvr_rep, err := c.Client.Do(req)
	if srvr_rep != nil {
		defer srvr_rep.Body.Close()
	}
	if err != nil {
		return nil, HttpError{err}
	}

	status := srvr_rep.StatusCode
	if status != http.StatusOK && status != http.StatusNotModified {
		return nil, newStatusError(srvr_rep)
	}

	cache_rep.Expires, err = getExpires(srvr_rep)
	if err != nil {
		return nil, MalformedResponse{err}
	}

	if status != http.StatusNotModified {
		err = json.NewDecoder(srvr_rep.Body).Decode(&cache_rep.Json)
		if err != nil {
			return nil, MalformedResponse{fmt.Errorf(
				"error decoding response body as json: %w",
				err,
			)}
		}
		cache_rep.Etag, err = getEtag(srvr_rep)
		if err != nil {
			return nil, MalformedResponse{err}
		}
	}

	err = c.Cache.Set(ctx, url, cache_rep)
	if err != nil {
		return cache_rep, CacheSetError{err} // rep is correct, but has not been placed into the cache
	}

	return cache_rep, nil
}

func (c *CrudeClient) RequestNoArray(
	ctx context.Context,
	url string,
	method string,
	auth string,
) (*response.EsiResponse, error) {
	c.Cache.Lock(url)
	defer c.Cache.Unlock(url)

	cache_rep, err := c.Cache.Get(ctx, url)
	if err != nil {
		return nil, CacheGetError{err}
	}
	if cache_rep != nil {
		now := time.Now()
		if now.Before(cache_rep.Expires) {
			return cache_rep, nil
		}
	} else {
		cache_rep = &response.EsiResponse{}
	}

	fmt.Printf("Sending '%s' request to '%s'\n", method, url)

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, RequestParamsError{err}
	}

	withHeadUserAgent(req)
	withHeadJsonContentType(req)
	withHeadAuthorization(req, auth)
	withHeadEtag(req, cache_rep.Etag)

	srvr_rep, err := c.Client.Do(req)
	if srvr_rep != nil {
		defer srvr_rep.Body.Close()
	}
	if err != nil {
		return nil, HttpError{err}
	}

	status := srvr_rep.StatusCode
	if status != http.StatusOK && status != http.StatusNotModified {
		return nil, newStatusError(srvr_rep)
	}

	cache_rep.Expires, err = getExpires(srvr_rep)
	if err != nil {
		return nil, MalformedResponse{err}
	}

	if status != http.StatusNotModified {
		no_array_json := new(map[string]interface{})
		err = json.NewDecoder(srvr_rep.Body).Decode(no_array_json)
		if err != nil {
			return nil, MalformedResponse{fmt.Errorf(
				"error decoding response body as json: %w",
				err,
			)}
		}
		array_json := make([]map[string]interface{}, 1)
		array_json[0] = *no_array_json
		cache_rep.Json = array_json
		cache_rep.Etag, err = getEtag(srvr_rep)
		if err != nil {
			return nil, MalformedResponse{err}
		}
	}

	err = c.Cache.Set(ctx, url, cache_rep)
	if err != nil {
		return cache_rep, CacheSetError{err} // rep is correct, but has not been placed into the cache
	}

	return cache_rep, nil
}

func (c *CrudeClient) RequestHead(
	ctx context.Context,
	url string,
	auth string,
) (*response.EsiHeadResponse, error) {
	c.Cache.LockHead(url)
	defer c.Cache.UnlockHead(url)

	cache_rep, err := c.Cache.GetHead(ctx, url)
	if err != nil {
		return nil, CacheGetError{err}
	}
	if cache_rep != nil {
		now := time.Now()
		if now.Before(cache_rep.Expires) {
			return cache_rep, nil
		}
	} else {
		cache_rep = &response.EsiHeadResponse{}
	}

	fmt.Printf("Sending 'HEAD' request to '%s'\n", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return nil, RequestParamsError{err}
	}

	withHeadUserAgent(req)
	withHeadAuthorization(req, auth)

	srvr_rep, err := c.Client.Do(req)
	if srvr_rep != nil {
		defer srvr_rep.Body.Close()
	}
	if err != nil {
		return nil, HttpError{err}
	}

	// fmt.Printf("Response Headers: %v", srvr_rep.Header)

	status := srvr_rep.StatusCode
	if status != http.StatusOK {
		return nil, newStatusError(srvr_rep)
	}

	cache_rep.Pages, err = getPages(srvr_rep)
	if err != nil {
		return nil, MalformedResponse{err}
	}
	cache_rep.Expires, err = getExpires(srvr_rep)
	if err != nil {
		return nil, MalformedResponse{err}
	}
	err = c.Cache.SetHead(ctx, url, cache_rep)
	if err != nil {
		return cache_rep, CacheSetError{err} // rep is correct, but has not been placed into the cache
	}

	return cache_rep, nil
}

type authenticationResponse struct {
	AccessToken string `json:"access_token"`
}

func (c *CrudeClient) RequestAuth(
	ctx context.Context,
	token string,
) (string, error) {
	var err error

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		AUTH_URL,
		bytes.NewBuffer([]byte(fmt.Sprintf(
			`grant_type=refresh_token&refresh_token=%s`,
			url.QueryEscape(token),
		))),
	)
	if err != nil {
		return "", RequestParamsError{err}
	}

	withHeadUserAgent(req)
	withHeadWwwContentType(req)
	withHeadAuthAuthorization(req)
	withHeadLoginHost(req)

	srvr_rep, err := c.Client.Do(req)
	if err != nil {
		return "", HttpError{err}
	}
	defer srvr_rep.Body.Close()

	status := srvr_rep.StatusCode
	if status != http.StatusOK {
		return "", newStatusError(srvr_rep)
	}

	auth_struct := new(authenticationResponse)
	err = json.NewDecoder(srvr_rep.Body).Decode(&auth_struct)
	if err != nil {
		return "", MalformedResponse{fmt.Errorf(
			"error decoding response body as json: %w",
			err,
		)}
	}

	return auth_struct.AccessToken, nil
}
