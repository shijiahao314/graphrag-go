package api

import (
	"bufio"
	"fmt"
	"graphraggo/internal/global"
	"graphraggo/internal/kb"
	"net/http"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
)

type KBApi struct {
}

func (ka *KBApi) Register(rg *gin.RouterGroup) {
	r := rg.Group("/kb")

	r.GET("", ka.GetKB)
	r.POST("indexing", ka.IndexingKB)
}

// GetKB 获取可用知识库
func (ka *KBApi) GetKB(c *gin.Context) {
	type GetKBReq struct {
	}
	type GetKBRsp struct {
		BaseRsp
		KBs []*kb.KB `json:"kbs"`
	}

	rsp := GetKBRsp{}

	kbs, err := kb.ReadKB()
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

// IndexingKB 建立索引
func (ka *KBApi) IndexingKB(c *gin.Context) {
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

	// cmd := exec.CommandContext(c, global.PythonPath, "--version")
	cmd := exec.CommandContext(c, global.PythonPath,
		"-m", "graphrag.index",
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
	c.Header("Content-Type", "text/plain")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)

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
