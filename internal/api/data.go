package api

import (
	"fmt"
	"graphraggo/internal/global"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type DataApi struct {
}

func (da *DataApi) Register(rg *gin.RouterGroup) {
	r := rg.Group("/data")

	r.POST("", da.GetData)
	r.POST("/delete", da.DeleteData)
}

// DeleteData 删除知识库
func (da *DataApi) DeleteData(c *gin.Context) {
	type DeleteDataReq struct {
		KB   string `json:"kb"`
		Name string `json:"name"`
	}
	type DeleteDataRsp struct {
		BaseRsp
	}

	req := DeleteDataReq{}
	rsp := DeleteDataRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// 判断文件夹是否存在
	path := fmt.Sprintf("%s/%s/%s/output/%s",
		global.WorkDir, global.KBDir, req.KB, req.Name)
	_, err := os.Stat(path)
	if err != nil {
		// 不存在
		rsp.Code = -1
		rsp.Msg = fmt.Sprintf("data '%s' not exists", req.Name)
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	if err := os.RemoveAll(path); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	c.JSON(http.StatusOK, rsp)
}

// ReadData 获取所有 Data
func ReadData(kb string) ([]string, error) {
	path := fmt.Sprintf("%s/%s/%s/output", global.WorkDir, global.KBDir, kb)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	kbs := []string{}
	for _, file := range files {
		if file.Type().IsDir() {
			kbs = append(kbs, file.Name())
		}
	}

	return kbs, nil
}

// GetData 获取可用 Data
func (da *DataApi) GetData(c *gin.Context) {
	type GetDataReq struct {
		KB string `json:"kb"`
	}
	type GetDataRsp struct {
		BaseRsp
		Datas []string `json:"datas"`
	}

	req := GetDataReq{}
	rsp := GetDataRsp{}

	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	datas, err := ReadData(req.KB)
	if err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	rsp.Datas = datas
	c.JSON(http.StatusOK, rsp)
}
