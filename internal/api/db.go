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
	r := rg.Group("/db")

	r.POST("", da.GetData)
	r.POST("/output", da.GetOutput)
	r.POST("/delete", da.DeleteData)
	r.POST("/logs", da.GetLogs)
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
	path := fmt.Sprintf("%s/%s/%s/output",
		global.WorkDir, global.KBDir, req.KB)
	_, err := os.Stat(path)
	if err != nil {
		// 不存在
		rsp.Code = -1
		rsp.Msg = fmt.Sprintf("db '%s' not exists", req.Name)
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
	return []string{"output"}, nil
}

// GetData 获取可用 Data
func (da *DataApi) GetData(c *gin.Context) {
	type GetDataReq struct {
		KB string `json:"kb"`
	}
	type GetDataRsp struct {
		BaseRsp
		DBs []string `json:"dbs"`
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
	rsp.DBs = datas
	c.JSON(http.StatusOK, rsp)
}

// ReadOutput 获取所有 Output
func ReadOutput(kb, db string) ([]string, error) {
	path := fmt.Sprintf("%s/%s/%s/output",
		global.WorkDir, global.KBDir, kb)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	outputs := []string{}
	for _, file := range files {
		if file.Type().IsDir() {
			continue
		}
		outputs = append(outputs, file.Name())
	}

	return outputs, nil
}

// GetOutput 获取输出文件
func (da *DataApi) GetOutput(c *gin.Context) {
	type GetDataReq struct {
		KB string `json:"kb"`
		DB string `json:"db"`
	}
	type GetDataRsp struct {
		BaseRsp
		Files []string `json:"files"`
	}

	req := GetDataReq{}
	rsp := GetDataRsp{}

	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	files, err := ReadOutput(req.KB, req.DB)
	if err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	rsp.Files = files
	c.JSON(http.StatusOK, rsp)
}

// ReadLogs 获取日志文件内容
func ReadLogs(kb, db string) ([]byte, error) {
	filename := "indexing-engine.log"
	logFilePath := fmt.Sprintf("%s/%s/%s/logs/%s",
		global.WorkDir, global.KBDir, kb, filename)

	files, err := os.ReadFile(logFilePath)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetLogs 获取日志文件
func (da *DataApi) GetLogs(c *gin.Context) {
	type GetDataReq struct {
		KB string `json:"kb"`
		DB string `json:"db"`
	}
	type GetDataRsp struct {
		BaseRsp
		Files string `json:"files"`
	}

	req := GetDataReq{}
	rsp := GetDataRsp{}

	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	files, err := ReadLogs(req.KB, req.DB)
	if err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	rsp.Files = string(files)
	c.JSON(http.StatusOK, rsp)
}
