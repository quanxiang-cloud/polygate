package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/polygate/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polygate/pkg/lib/httputil"
)

// PingPong return a pong response for ping
func PingPong(c *gin.Context) {
	ret, err := ping(c)
	resp.Format(ret, err).Context(c)
}

func ping(c *gin.Context) (*polysign.PingResp, error) {
	req := polysign.PingReq{}
	if err := httputil.BindBody(c, &req); err != nil {
		return nil, err
	}
	ret := &polysign.PingResp{
		PingReq:       req,
		PongTimestamp: time.Now().Format(polysign.PingTimestampFmt),
	}
	return ret, nil
}
