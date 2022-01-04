package gateentry

import (
	"net/http"

	"github.com/quanxiang-cloud/cabin/logger"
	ginlog "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/polygate/pkg/config"
	"github.com/quanxiang-cloud/polygate/pkg/gate/chain"
	"github.com/quanxiang-cloud/polygate/pkg/gate/chain/regist"
	"github.com/quanxiang-cloud/polygate/pkg/gate/impl/authorize"
	"github.com/quanxiang-cloud/polygate/pkg/gate/impl/ginproxy"
	"github.com/quanxiang-cloud/polygate/pkg/gate/impl/goalie"
	"github.com/quanxiang-cloud/polygate/pkg/gate/impl/requestid"

	"github.com/gin-gonic/gin"
)

// InitGate init regist of gate nodes
func InitGate(cfg *config.Config) error {
	regist.MustRegist("", authorize.New(cfg))
	regist.MustRegist("", goalie.New(cfg))
	regist.MustRegist("", requestid.New(cfg))
	regist.MustRegist("", ginproxy.New(cfg))
	return nil
}

// New create a gate entry
func New(cfg *config.Config) (*GateEntry, error) {
	return &GateEntry{
		n: regist.ToChain(),
	}, nil
}

//------------------------------------------------------------------------------

// GateEntry is entry of gate
type GateEntry struct {
	n *chain.Node
}

// Handle is the main handler of gate
func (v *GateEntry) Handle(c *gin.Context) {
	if err := v.n.Handle(c); err != nil {
		logger.Logger.PutError(err, "GateEntry", ginlog.GINRequestID(c))
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
	}
}
