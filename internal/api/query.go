package api

import (
	"fmt"
	"graphraggo/internal/global"
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

type Method string

const (
	Local  Method = "local"
	Global Method = "global"
)

func (qa *QueryApi) Query(c *gin.Context) {
	type QueryReq struct {
		KB        string `json:"kb"`
		Timestamp string `json:"timestamp"`
		Method    Method `json:"method"`
		Text      string `json:"text"`
	}
	type QueryRsp struct {
		BaseRsp
		Text string `jons:"text"`
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
	config := fmt.Sprintf("%s/%s", path, "settings.yaml")
	data := fmt.Sprintf("%s/output/%s/artifacts", path, req.Timestamp)

	cmd := exec.CommandContext(c, global.PythonPath,
		"-m", "graphrag.query",
		"--config", config,
		"--data", data,
		"--method", string(req.Method),
		"--response_type", "Single Paragraph",
		req.Text)

	out, err := cmd.CombinedOutput()
	if err != nil {
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
