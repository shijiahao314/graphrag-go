# Knowledge Base

本文件夹存储知识库，每个知识库为一个文件夹

本文件夹默认自带一个官方实例知识库 `ragtest` 和一个修改了配置的知识库 `raggo`

## 激活 conda 环境

```bash
conda activate graphrag-go
graphrag init --root ./ragtest
```

## 建立知识库索引（Index）

```bash
graphrag index --root ./raggo
```

## 问答（Query）

### Local

```bash
graphrag query \
--root ./raggo \
--method local \
--query "Who is Scrooge, and what are his main relationships?"
```

### Global

```bash
graphrag query \
--root ./raggo \
--method global \
--query "What are the top themes in this story?"
```

## 常见问题

### ValueError: Columns must be same length as key

修改 `settings.yaml` 中的 `chunks` 大小：

```yaml
chunks:
  size: 512
  overlap: 64
```
