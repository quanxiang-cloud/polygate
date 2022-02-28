package authorize

import (
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/polygate/pkg/config"
	"github.com/quanxiang-cloud/polygate/pkg/dependence/remotes"
	"github.com/quanxiang-cloud/polygate/pkg/gate/chain"
	"github.com/quanxiang-cloud/polygate/pkg/lib/httputil"
	"github.com/quanxiang-cloud/polygate/pkg/lib/tiretree"

	"github.com/gin-gonic/gin"
)

// New create a chain handler
func New(cfg *config.Config) chain.Handler {
	authKey := mustNewClient(&cfg.Remotes.OauthKey)
	authToken := mustNewClient(&cfg.Remotes.OauthToken)
	tt := tiretree.NewTireTree()
	if err := tt.BatchInsert(cfg.APIFilterConfig.White, tiretree.White); err != nil {
		panic(err)
	}
	logger.Logger.Infow(tt.Show())

	return &authorize{
		c:      remotes.NewAuthClient(authKey, authToken),
		filter: tt,
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

type authorize struct {
	c      remotes.Client
	filter *tiretree.TireTree
}

func (v *authorize) Handle(c *gin.Context) error {
	if b, ok := v.filter.Match(c.Request.URL.Path); ok && b == tiretree.White {
		return nil
	}
	return v.c.Request(c)
}

func (v *authorize) Name() string {
	return "authorize"
}
