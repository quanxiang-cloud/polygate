package ginproxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/quanxiang-cloud/polygate/pkg/basic/consts"
	"github.com/quanxiang-cloud/polygate/pkg/basic/errcode"
	"github.com/quanxiang-cloud/polygate/pkg/basic/header"
	"github.com/quanxiang-cloud/polygate/pkg/config"
	"github.com/quanxiang-cloud/polygate/pkg/gate/chain"

	"github.com/gin-gonic/gin"
)

// New create a chain handler
func New(cfg *config.Config) chain.Handler {
	proxyCfg := cfg.Proxy
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   proxyCfg.Timeout * time.Second,
			KeepAlive: proxyCfg.KeepAlive * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          proxyCfg.MaxIdleConns,
		IdleConnTimeout:       proxyCfg.IdleConnTimeout * time.Second,
		TLSHandshakeTimeout:   proxyCfg.TLSHandshakeTimeout * time.Second,
		ExpectContinueTimeout: proxyCfg.ExpectContinueTimeout * time.Second,
	}
	return &ginproxy{
		transport: transport,
		schema:    cfg.Schema,
	}
}

//------------------------------------------------------------------------------

type ginproxy struct {
	schema    string
	transport *http.Transport
}

func (v *ginproxy) Handle(c *gin.Context) error {
	url, err := v.reWriteURL(c)
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = v.transport
	r := c.Request
	r.Host = url.Host
	proxy.ServeHTTP(c.Writer, r)

	return nil
}

func (v *ginproxy) reWriteURL(c *gin.Context) (*url.URL, error) {
	dnsName := ""
	if s, ok := c.Get(header.HeaderXRedirectService); ok {
		if ss, ok := s.(string); ok {
			dnsName = ss
		}
	}
	if dnsName == "" {
		if s, ok := c.Params.Get(consts.PathArgServiceName); ok {
			dnsName = s
		} else {
			return nil, errcode.ErrInvalidURI.FmtError(c.Request.URL.String())
		}
	}

	return url.ParseRequestURI(v.schema + dnsName)
}

func (v *ginproxy) Name() string {
	return "ginproxy"
}
