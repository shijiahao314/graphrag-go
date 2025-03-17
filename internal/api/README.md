# Api

## NER

```bash
curl -X POST localhost:8080/api/ner \
  -H "Content-Type: application/json" \
  -d '{"text": "2021年HanLPv2.1为生产环境带来次世代最先进的多语种NLP技术。阿婆主来到北京立方庭参观自然语义科技公司。"}'
```

```bash
curl 'http://127.0.0.1:8081/ner' \
  -H 'Content-Type: application/json' \
  --data-raw '{"text":"萨哈夫说，伊拉克将同联合国销毁伊拉克大规模杀伤性武器特别委员会继续保持合作。"}'
```

## kb

### add

```bash
curl -X POST localhost:8080/api/kb/add \
  -H "Content-Type: application/json" \
  -d '{"name": "santi"}'
```

### delete

```bash
curl -X POST localhost:8080/api/kb/delete \
  -H "Content-Type: application/json" \
  -d '{"name": "santi"}'
```

### get

```bash
curl localhost:8080/api/kb
```

### indexing

```bash
curl -X POST localhost:8080/api/kb/indexing \
  -H "Content-Type: application/json" \
  -d '{"name": "raggo"}'
```

## db

### get

```bash
curl -X POST localhost:8080/api/db \
  -H "Content-Type: application/json" \
  -d '{"kb": "raggo"}'
```

### delete

```bash
# dont use it easily
curl -X POST localhost:8080/api/db/delete \
  -H "Content-Type: application/json" \
  -d '{"kb": "raggo", "name": "yyyyMMdd-hhmmss"}'
```

### logs

```bash
curl -X POST localhost:8080/api/db/logs \
  -H "Content-Type: application/json" \
  -d '{"kb": "raggo", "db": "yyyyMMdd-hhmmss"}'
```

## query

### local

需要使用绝对路径

```bash
python -m graphrag.query \
--config $(pwd)/kb/raggo/settings.yaml \
--data $(pwd)/kb/raggo/output/yyyyMMdd-hhmmss/artifacts \
--method local \
--response_type "Single Paragraph" \
"Who is Scrooge and what are his main relationships?"
```

```bash
curl -X POST localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{"kb": "raggo", "timestamp": "yyyyMMdd-hhmmss", "method": "local", "text": "Who is Scrooge and what are his main relationships?"}'
```

### global

```bash
python -m graphrag.query \
--config $(pwd)/kb/raggo/settings.yaml \
--data $(pwd)/kb/raggo/output/yyyyMMdd-hhmmss/artifacts \
--method global \
--response_type "Single Paragraph" \
"What are the top themes in this story?"
```

```bash
curl -X POST localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{"kb": "raggo", "timestamp": "yyyyMMdd-hhmmss", "method": "global", "text": "Who is Scrooge, and what are his main relationships?"}'
```
