package router

import (
	"fmt"

	"github.com/quanxiang-cloud/polygate/pkg/basic/consts"
	"github.com/quanxiang-cloud/polygate/pkg/config"
	"github.com/quanxiang-cloud/polygate/pkg/gate/gateentry"
	"github.com/quanxiang-cloud/polygate/pkg/lib/ginlog"
	"github.com/quanxiang-cloud/polygate/pkg/probe"

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
	*probe.Probe

	c *config.Config

	engine *gin.Engine
}

// NewRouter create router
func NewRouter(c *config.Config) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}
	r := &Router{
		c:      c,
		engine: engine,
		Probe:  probe.New(),
	}

	r.probe()
	return r, nil
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

	apiPath := fmt.Sprintf("/api/:version/:%s/*realPath", consts.PathArgServiceName)
	engine.Any(apiPath, entry.Handle)
	engine.POST("/api/v1/gate/ping", PingPong)

	return engine, nil
}

func (r *Router) probe() {
	r.engine.GET("liveness", func(c *gin.Context) {
		r.Probe.LivenessProbe(c.Writer, c.Request)
	})

	r.engine.Any("readiness", func(c *gin.Context) {
		r.Probe.ReadinessProbe(c.Writer, c.Request)
	})
}

// Run start router
func (r *Router) Run() {
	r.Probe.SetRunning()
	r.engine.Run(r.c.Port)
}

// Close close router
func (r *Router) Close() {
}
