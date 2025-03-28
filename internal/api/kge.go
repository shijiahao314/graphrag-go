package api

import (
	"fmt"
	"graphraggo/internal/global"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	imageFolder = "image"
)

type KGEApi struct {
}

func (kgeApi *KGEApi) Register(rg *gin.RouterGroup) {
	r := rg.Group("/kge")

	r.POST("", kgeApi.KGEGraph)
}

type KGEReq struct {
	Image string `json:"image"`
}

type KGERsp struct {
	BaseRsp
}

func (kgeApi *KGEApi) KGEGraph(c *gin.Context) {
	req := KGEReq{}
	rsp := KGERsp{}

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("request invalid", slog.String("error", err.Error()))
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// 获取本地图片
	imagePath := fmt.Sprintf("%s/py/%s/%s.png", global.WorkDir, imageFolder, req.Image)
	fmt.Println(imagePath)

	// 确认图片存在并返回
	if _, err := os.Stat(imagePath); err == nil {
		// 返回图片文件
		c.File(imagePath)
	} else {
		rsp.Code = -1
		rsp.Msg = "Image not found"
		c.JSON(http.StatusInternalServerError, rsp)
	}
}
