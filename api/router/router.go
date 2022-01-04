package router

import (
	"fmt"

	ginlog "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/polygate/pkg/basic/consts"
	"github.com/quanxiang-cloud/polygate/pkg/config"
	"github.com/quanxiang-cloud/polygate/pkg/gate/gateentry"

	"github.com/gin-gonic/gin"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router router
type Router struct {
	c *config.Config

	engine *gin.Engine
}

// NewRouter create router
func NewRouter(c *config.Config) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}
	return &Router{
		c:      c,
		engine: engine,
	}, nil
}

func newRouter(c *config.Config) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()

	engine.Use(ginlog.GinLogger(), ginlog.GinRecovery())
	if err := gateentry.InitGate(c); err != nil {
		return nil, err
	}
	entry, err := gateentry.New(c)
	if err != nil {
		return nil, err
	}

	apiPath := fmt.Sprintf("/api/v1/:%s/*realPath", consts.PathArgServiceName)
	engine.Any(apiPath, entry.Handle)

	return engine, nil
}

// Run start router
func (r *Router) Run() {
	r.engine.Run(r.c.Port)
}

// Close close router
func (r *Router) Close() {
}
