package api

import (
	"bufio"
	"fmt"
	"graphraggo/internal/global"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
)

type KBApi struct {
}

func (ka *KBApi) Register(rg *gin.RouterGroup) {
	r := rg.Group("/kb")

	r.GET("", ka.GetKB)
	r.POST("/input", ka.GetInput)
	r.POST("/add", ka.AddKB)
	r.POST("/delete", ka.DeleteKB)
	r.POST("/indexing", ka.IndexKB)
}

// AddKB 新建知识库
func (ka *KBApi) AddKB(c *gin.Context) {
	type AddKBReq struct {
		Name string `json:"name"`
	}
	type AddKBRsp struct {
		BaseRsp
	}

	req := AddKBReq{}
	rsp := AddKBRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// 判断文件夹是否存在
	path := fmt.Sprintf("%s/%s/%s", global.WorkDir, global.KBDir, req.Name)
	_, err := os.Stat(path)
	if err == nil {
		// 存在
		rsp.Code = -1
		rsp.Msg = fmt.Sprintf("kb '%s' already exists", req.Name)
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	if !os.IsNotExist(err) {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	os.Mkdir(path, os.ModePerm)
	os.Mkdir(path+"/input", os.ModePerm)
	cmd := exec.Command("cp", global.ExampleSettingFile, path+"/settings.yaml")
	if err := cmd.Run(); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	c.JSON(http.StatusOK, rsp)
}

// DeleteKB 删除知识库
func (ka *KBApi) DeleteKB(c *gin.Context) {
	type DeleteKBReq struct {
		Name string `json:"name"`
	}
	type DeleteKBRsp struct {
		BaseRsp
	}

	req := DeleteKBReq{}
	rsp := DeleteKBRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	if req.Name == "" {
		rsp.Code = -1
		rsp.Msg = "kb name is empty"
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// 判断文件夹是否存在
	path := fmt.Sprintf("%s/%s/%s", global.WorkDir, global.KBDir, req.Name)
	_, err := os.Stat(path)
	if err != nil {
		// 不存在
		rsp.Code = -1
		rsp.Msg = fmt.Sprintf("kb '%s' not exists", req.Name)
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

// ReadKB 获取所有知识库
func ReadKB() ([]string, error) {
	path := fmt.Sprintf("%s/%s", global.WorkDir, global.KBDir)

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

// GetKB 获取可用知识库
func (ka *KBApi) GetKB(c *gin.Context) {
	type GetKBReq struct {
	}
	type GetKBRsp struct {
		BaseRsp
		KBs []string `json:"kbs"`
	}

	rsp := GetKBRsp{}

	kbs, err := ReadKB()
	if err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	rsp.KBs = kbs
	c.JSON(http.StatusOK, rsp)
}

var lock sync.Mutex

// IndexKB 建立索引
func (ka *KBApi) IndexKB(c *gin.Context) {
	type IndexingKBReq struct {
		Name string `json:"name"`
	}
	type IndexingKBRsp struct {
		BaseRsp
	}

	req := IndexingKBReq{}
	rsp := IndexingKBRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	path := fmt.Sprintf("%s/%s/%s", global.WorkDir, global.KBDir, req.Name)

	// 创建文件记录状态
	if !lock.TryLock() {
		// 加锁失败
		rsp.Code = -1
		rsp.Msg = "already in indexing process"
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	defer lock.Unlock() // 释放锁

	cmd := exec.CommandContext(c, global.PythonPath,
		"-m", "graphrag index",
		"--root", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	if err := cmd.Start(); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	// 设置流式响应的头部
	// c.Header("Content-Type", "text/plain")
	// c.Header("Transfer-Encoding", "chunked")
	// c.Status(http.StatusOK)

	// 创建一个 bufio 读取器用于逐行读取命令输出
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		c.Writer.Write([]byte(line)) // 将读取到的行写入到客户端
		c.Writer.Flush()             // 刷新缓冲区，确保数据立刻发送
	}

	// 等待命令执行完毕
	if err := cmd.Wait(); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Code = 0
	rsp.Msg = "success"
	c.JSON(http.StatusOK, rsp)
}

// ReadInput 获取所有 Input
func ReadInput(kb string) ([]string, error) {
	path := fmt.Sprintf("%s/%s/%s/input",
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

// GetInput 获取输入文件
func (da *KBApi) GetInput(c *gin.Context) {
	type GetDataReq struct {
		KB string `json:"kb"`
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

	files, err := ReadInput(req.KB)
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
