# GraphRAG-Go

使用 Go 提供 GraphRAG API 服务

## 1.安装 GraphRAG（0.5.0） 环境

官方网站：[GraphRAG](https://microsoft.github.io/graphrag/)

### 创建 conda 环境

```bash
# conda env remove -n graphrag-go
conda create -n graphrag-go python=3.12 -y
conda activate graphrag-go
pip install graphrag==0.5.0
pip install ollama
# python ner server
pip install hanlp, fastapi, uvicorn
```

## 2.修改本地 conda 环境 graphrrag 包源代码

参考 [change/README.md](./change/README.md)

## 3.测试

参考 [kb/README.md](./kb/README.md)

## 4.运行服务

```bash
go run main.go
```

## 5.测试 API

参考 [internal/api/README.md](./internal/api/README.md)
