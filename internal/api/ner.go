package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"graphraggo/internal/global"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NERApi struct {
}

func (na *NERApi) Register(rg *gin.RouterGroup) {
	r := rg.Group("/ner")

	r.POST("", na.NER)
}

type NERReq struct {
	Text string `json:"text"`
}

type NERRsp struct {
	BaseRsp
	Text string `json:"text"`
}

func (na *NERApi) NER(c *gin.Context) {

	req := NERReq{}
	rsp := NERRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// 调用 NER 服务
	result, err := callNERService(req.Text)
	if err != nil {
		slog.Error("decode error", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, result)
}

func callNERService(text string) (*NERRsp, error) {
	url := fmt.Sprintf("http://%s:%d/ner", global.Host, global.PythonServerPort)

	// 构造请求体
	reqBody, err := json.Marshal(NERReq{Text: text})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 创建 HTTP 客户端，设置超时时间
	client := &http.Client{Timeout: 5 * time.Second}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查返回状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// 解析 JSON 响应
	var result NERRsp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil
}
