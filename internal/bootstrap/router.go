package bootstrap

import (
	"graphraggo/internal/api"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	HealthPath = "/api/health"
)

type IRouter interface {
	Register(r *gin.RouterGroup)
}

// MustInitRouter 初始化路由配置
func MustInitRouter() *gin.Engine {
	r := gin.New()

	r.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{HealthPath}}),
		gin.Recovery(),
	)

	r.Use(cors.Default())

	r.GET(HealthPath, func(c *gin.Context) {
		c.String(http.StatusOK, "ok!")
	})

	g := r.Group("/api")

	routers := []IRouter{
		&api.NERApi{},
		&api.KBApi{},
		&api.DataApi{},
		&api.QueryApi{},
	}
	for _, rt := range routers {
		rt.Register(g)
	}

	return r
}
