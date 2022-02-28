package httputil

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

// HTTPClientConfig http client
type HTTPClientConfig struct {
	Addr              string
	MaxIdleConns      int
	TimeoutSec        time.Duration
	DisableKeepAlives bool
}

// MakeHTTPClientConfig create a http client config
func MakeHTTPClientConfig(addr string, timeoutSec time.Duration, maxIdleConn int, keepAlive bool) *HTTPClientConfig {
	return &HTTPClientConfig{
		Addr:              addr,
		TimeoutSec:        timeoutSec,
		MaxIdleConns:      maxIdleConn,
		DisableKeepAlives: !keepAlive,
	}
}

// HTTPClient object
type HTTPClient struct {
	URI    *url.URL
	Client http.Client
}

// NewHTTPClient create a reuseable http client
func NewHTTPClient(cfg *HTTPClientConfig) (*HTTPClient, error) {
	uri, err := url.ParseRequestURI(cfg.Addr)
	if err != nil {
		return nil, err
	}

	c := &HTTPClient{
		URI: uri,
		Client: http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					timeout := cfg.TimeoutSec * time.Second
					deadline := time.Now().Add(timeout)
					c, err := net.DialTimeout(netw, addr, timeout)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
				MaxIdleConns:      cfg.MaxIdleConns,
				DisableKeepAlives: cfg.DisableKeepAlives,
			},
		},
	}
	return c, nil
}
