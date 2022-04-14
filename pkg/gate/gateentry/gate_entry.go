package gateentry

import (
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/polygate/pkg/basic/errcode"
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
	regist.MustRegist("", requestid.New(cfg))
	regist.MustRegist("", authorize.New(cfg))
	regist.MustRegist("", goalie.New(cfg))
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
	var err error
	for p := v.n; p != nil; p = p.Next {
		if e := p.H.Handle(c); e != nil {
			err = errcode.ErrGateError.FmtError(p.GetName(), e.Error())
			break
		}
	}
	if err != nil {
		resp.Format(nil, err).Context(c, c.Writer.Status())
	}
}
