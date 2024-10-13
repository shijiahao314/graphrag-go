# GraphRAG-Go

使用 Go 提供 GraphRAG + Ollama 后端服务

## 安装 GraphRAG 环境

https://microsoft.github.io/graphrag/

所需：

- conda
- ollama

### 创建 conda 环境

```bash
conda create -n graphrag-go python=3.12
conda activate graphrag-go
pip install graphrag
pip install ollama
```

### 运行 ollama

```bash
# 确保 ollama 处于服务中
ollama -v # ollama version is 0.3.12
# 拉取 llm
ollama pull $LLM
# 拉取 embedding
ollama pull $EMBEDDING
```
