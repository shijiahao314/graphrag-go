package api

import (
	"fmt"
	"graphraggo/internal/global"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type QueryApi struct {
}

func (qa *QueryApi) Register(rg *gin.RouterGroup) {
	r := rg.Group("/query")

	r.POST("", qa.Query)
}

type QueryMethod string

const (
	Local  QueryMethod = "local"
	Global QueryMethod = "global"
)

// Query 提供查询能力
func (qa *QueryApi) Query(c *gin.Context) {
	type QueryReq struct {
		KB     string      `json:"kb"`
		DB     string      `json:"db"`
		Method QueryMethod `json:"method"`
		Text   string      `json:"text"`
	}
	type QueryRsp struct {
		BaseRsp
		Text string `json:"text"`
	}

	req := QueryReq{}
	rsp := QueryRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	path := fmt.Sprintf("%s/%s/%s", global.WorkDir, global.KBDir, req.KB)
	query := strings.Replace(req.Text, "\n", "\\n", -1)

	cmd := exec.CommandContext(c, global.PythonPath,
		"-m", "graphrag", "query",
		"--root", path,
		"--method", string(req.Method),
		"--query", query, // 使用转义后的查询文本
		"--response-type", "Single Paragraph",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error(err.Error(), slog.String("cmd", cmd.String()))
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	info := string(out)
	search := "Search Response:"
	n := strings.Index(info, search)
	if n != -1 {
		info = info[n+len(search)+1:]
	}
	rsp.Text = info
	c.JSON(http.StatusOK, rsp)
}
