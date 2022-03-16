package remotes

import (
	"encoding/json"
	"net/http"

	"github.com/quanxiang-cloud/cabin/logger"
	ginlog "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/polygate/pkg/basic/header"
	"github.com/quanxiang-cloud/polygate/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polygate/pkg/config"
	"github.com/quanxiang-cloud/polygate/pkg/lib/httputil"

	"github.com/gin-gonic/gin"
)

var ginRequestID = ginlog.GINRequestID // function exports
var log = logger.Logger.WithName("dependence.remotes")

// Client is a remote serivice client
type Client interface {
	Request(c *gin.Context) error
}

// NewHTTPClient create a http client by config
func NewHTTPClient(cfg *config.HTTPClientConfig) (*httputil.HTTPClient, error) {
	arg := httputil.MakeHTTPClientConfig(cfg.Addr, cfg.Timeout, cfg.MaxIdleConns, true)
	return httputil.NewHTTPClient(arg)
}

//------------------------------------------------------------------------------

// RequestArg is arg of request
type requestArg struct {
	signature   string
	accessKeyID string
	body        json.RawMessage
}

type resolver interface {
	request(c *gin.Context, arg *requestArg) error
}

func cloneProfile(dst *http.Header, src http.Header) {
	dst.Set(header.HeaderUserID, header.DeepCopy(src.Values(header.HeaderUserID)))
	dst.Set(header.HeaderUserName, header.DeepCopy(src.Values(header.HeaderUserName)))
	dst.Set(header.HeaderDepartmentID, header.DeepCopy(src.Values(header.HeaderDepartmentID)))
	dst.Set(header.HeaderTenantID, header.DeepCopy(src.Values(header.HeaderTenantID)))
}

func removeAuthArgs(c *gin.Context) {
	delete(c.Request.Header, polysign.XHeaderPolySignKeyID)
	delete(c.Request.Header, polysign.XHeaderPolySignMethod)
	delete(c.Request.Header, polysign.XHeaderPolySignVersion)
	delete(c.Request.Header, polysign.XHeaderPolySignTimestamp)
	delete(c.Request.Header, header.HeaderAccessToken)
}
