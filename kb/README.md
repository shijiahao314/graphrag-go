# Knowledge Base

本文件夹存储知识库，每个知识库为一个文件夹

本文件夹默认自带一个官方实例知识库 `ragtest` 和一个修改了配置的知识库 `raggo`

## 运行

```bash
conda activate graphrag-go
python -m graphrag.index --root ./raggo
```

## 问答

Local:

```bash
python -m graphrag.query \
--root ./raggo \
--method local \
"Who is Scrooge, and what are his main relationships?"
```

Global:

```bash
python -m graphrag.query \
--root ./raggo \
--method global \
"What are the top themes in this story?"
```