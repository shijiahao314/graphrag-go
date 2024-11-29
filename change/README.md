# Change

Python 安装包 graphrag 中需要修改的文件

## 所需替换的文件

- graphrag/llm/openai/openai_embeddings_llm.py
- graphrag/query/llm/oai/embedding.py
- graphrag/query/llm/text_utils.py

## 查看 graphrag 路径

激活所使用的 conda 环境

```bash
conda activate <env_name>
```

查看 graphrag 路径位置

```bash
pip show graphrag | grep Location
```

输出结果类似

```text
Location: /home/shijiahao/miniconda3/envs/graphrag-go/lib/python3.12/site-packages
```

替换文件（确保在本文件所在目录执行）

```bash
# 设置环境变量
export PACKAGE_PATH="/home/shijiahao/miniconda3/envs/graphrag-go/lib/python3.12/site-packages"
# 替换文件
cp ./openai_embeddings_llm.py $PACKAGE_PATH/graphrag/llm/openai/openai_embeddings_llm.py
cp ./embedding.py $PACKAGE_PATH/graphrag/query/llm/oai/embedding.py
# cp ./text_utils.py $PACKAGE_PATH/graphrag/query/llm/text_utils.py
# 解除环境变量
unset PACKAGE_PATH
```
