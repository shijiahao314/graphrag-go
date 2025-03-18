package bootstrap

import (
	"fmt"
	"graphraggo/internal/api"
	"graphraggo/internal/global"
	"log/slog"
	"net/http"
	"os/exec"
	"time"

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
		&api.KGCApi{},
		&api.KBApi{},
		&api.DataApi{},
		&api.QueryApi{},
	}
	for _, rt := range routers {
		rt.Register(g)
	}

	return r
}

// MustInitPythonServer 启动Python服务
func MustInitPythonServer() {
	nerServer := fmt.Sprintf("%s/%s", global.WorkDir, "/py/py_server.py")

	cmd := exec.Command(global.PythonPath, nerServer,
		"--host", global.Host,
		"--port", fmt.Sprint(global.PythonServerPort))
	if err := cmd.Start(); err != nil {
		slog.Error("failed to run ner_server",
			slog.String("host", global.Host),
			slog.Int("port", global.PythonServerPort),
			slog.String("err", err.Error()))
		return
	}

	// 等待服务启动
	time.Sleep(5 * time.Second)

	// 检查服务状态
	checkServerStatus := func() {
		url := fmt.Sprintf("http://%s:%d/docs", global.Host, global.PythonServerPort)

		// 尝试 10 次
		for i := 0; i < 10; i++ {
			rsp, err := exec.Command("curl", "-I", url).Output()
			if err != nil {
				slog.Error("failed to check server status",
					slog.String("err", err.Error()))
				time.Sleep(5 * time.Second)
				continue
			}

			slog.Info("server is ready",
				slog.String("rsp", string(rsp)))
			break
		}
	}

	checkServerStatus()

	select {}
}
