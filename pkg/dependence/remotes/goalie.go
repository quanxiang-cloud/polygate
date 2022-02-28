package remotes

import (
	"net/http"

	"github.com/quanxiang-cloud/polygate/pkg/basic/errcode"
	"github.com/quanxiang-cloud/polygate/pkg/basic/header"
	"github.com/quanxiang-cloud/polygate/pkg/lib/httputil"

	"github.com/gin-gonic/gin"
)

var _ Client = (*goalie)(nil)

// NewGoalieClient create a goalie client
func NewGoalieClient(client *httputil.HTTPClient) Client {
	return (*goalie)(client)
}

type goalie httputil.HTTPClient

func (r *goalie) Request(c *gin.Context) error {
	req := &http.Request{
		Header: c.Request.Header.Clone(),
		URL:    r.URI,
	}
	req.Header.Set(header.HeaderContentType, header.MIMEJSON)

	resp, err := r.Client.Do(req)
	if err != nil {
		log.PutError(err, "Goalie.request", ginRequestID(c))
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(resp.StatusCode)
		return errcode.ErrInternal.FmtError(resp.Status)
	}

	c.Request.Header.Set(header.HeaderRole, resp.Header.Get(header.HeaderRole))
	return nil
}
