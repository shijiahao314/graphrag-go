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

type KGCApi struct {
}

func (kgcApi *KGCApi) Register(rg *gin.RouterGroup) {
	r := rg.Group("/kgc")

	r.POST("", kgcApi.KGC)
	r.GET("/benchmark", kgcApi.KGCBenckmark)
}

type KGCReq struct {
	Head     string `json:"head"`
	Relation string `json:"relation"`
	Tail     string `json:"tail"`
}

type KGCRsp struct {
	BaseRsp
	Head     string `json:"head"`
	Relation string `json:"relation"`
	Tail     string `json:"tail"`
}

func (kgcApi *KGCApi) KGC(c *gin.Context) {
	req := KGCReq{}
	rsp := KGCRsp{}

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("request invalid", slog.String("error", err.Error()))
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// 调用 KGC 服务
	result, err := callKGCService(req.Head, req.Relation, req.Tail)
	if err != nil {
		slog.Error("call kgc service error", slog.String("error", err.Error()))
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	fmt.Println(result)

	// 返回结果
	c.JSON(http.StatusOK, result)
}

func callKGCService(head, relation, tail string) (*KGCRsp, error) {
	url := fmt.Sprintf("http://127.0.0.1:%d/kgc", global.PythonServerPort)

	// 构造请求体
	reqBody, err := json.Marshal(KGCReq{Head: head, Relation: relation, Tail: tail})
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
	var result KGCRsp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil
}

type KGCBenckmarkRsp struct {
	BaseRsp
	HitsAt1 string `json:"hits_at_1"`
	MRR     string `json:"mrr"`
}

func (kgcApi *KGCApi) KGCBenckmark(c *gin.Context) {
	rsp := KGCBenckmarkRsp{}

	// 调用 KGC 服务
	result, err := callKGCBenckmarkService()
	if err != nil {
		slog.Error("call kgc benchmark service error", slog.String("error", err.Error()))
		rsp.Code = -1
		rsp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	fmt.Println(result)

	// 返回结果
	c.JSON(http.StatusOK, result)
}

func callKGCBenckmarkService() (*KGCBenckmarkRsp, error) {
	url := fmt.Sprintf("http://127.0.0.1:%d/kgc_benchmark", global.PythonServerPort)

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
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
	var result KGCBenckmarkRsp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil
}
