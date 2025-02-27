package api

import (
	"fmt"
	"graphraggo/internal/global"
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type NERApi struct {
}

func (na *NERApi) Register(rg *gin.RouterGroup) {
	r := rg.Group("/ner")

	r.POST("", na.NER)
}

func (na *NERApi) NER(c *gin.Context) {
	type NERReq struct {
		Text string `json:"text"`
	}

	type NERRsp struct {
		BaseRsp
		Text string `json:"text"`
	}

	req := NERReq{}
	rsp := NERRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	pythonFilePath := fmt.Sprintf("%s/%s", global.WorkDir, "/py/ner.py")

	cmd := exec.CommandContext(c, global.PythonPath,
		pythonFilePath, req.Text)

	out, err := cmd.Output()
	if err != nil {
		slog.Error(err.Error(), slog.String("cmd", cmd.String()))
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	rsp.Text = string(out)
	c.JSON(http.StatusOK, rsp)
}
