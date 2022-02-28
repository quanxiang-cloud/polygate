package requestid

import (
	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/polygate/pkg/basic/header"
	"github.com/quanxiang-cloud/polygate/pkg/config"
	"github.com/quanxiang-cloud/polygate/pkg/gate/chain"

	"github.com/gin-gonic/gin"
)

// New create a chain handler
func New(cfg *config.Config) chain.Handler {
	return &requsetid{}
}

//------------------------------------------------------------------------------

type requsetid struct{}

func (v *requsetid) Handle(c *gin.Context) error {
	id := c.Request.Header.Get(header.HeaderXRequestID)
	if len(id) != 0 {
		id += "-"
	}
	c.Request.Header.Add(header.HeaderRequestID, id+v.genID())
	return nil
}

func (v *requsetid) genID() string {
	return id.WithPrefix(id.ShortID(12), "req_")
}

func (v *requsetid) Name() string {
	return "requsetid"
}
