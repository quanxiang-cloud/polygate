package goalie

import (
	"github.com/quanxiang-cloud/polygate/pkg/config"
	"github.com/quanxiang-cloud/polygate/pkg/dependence/remotes"
	"github.com/quanxiang-cloud/polygate/pkg/gate/chain"
	"github.com/quanxiang-cloud/polygate/pkg/lib/httputil"

	"github.com/gin-gonic/gin"
)

// New create a chain handler
func New(cfg *config.Config) chain.Handler {
	c := mustNewClient(&cfg.Remotes.Goalie)

	return &goalie{
		c: remotes.NewGoalieClient(c),
	}
}

//------------------------------------------------------------------------------

func mustNewClient(cfg *config.HTTPClientConfig) *httputil.HTTPClient {
	client, err := remotes.NewHTTPClient(cfg)
	if err != nil {
		panic(err)
	}
	return client
}

type goalie struct {
	c remotes.Client
}

func (v *goalie) Handle(c *gin.Context) error {
	return v.c.Request(c)
}

func (v *goalie) Name() string {
	return "goalie"
}
