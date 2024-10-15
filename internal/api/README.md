# Api

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

## query

### ask

local: 

```bash
python -m graphrag.query \
--config ./kb/raggo/settings.yaml \
--data ./kb/raggo/output/20241015-132435/artifacts \
--method local \
--response_type "Single Paragraph" \
"What are the top themes in this story?"
```

```bash
curl -X POST localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{"kb": "raggo", "timestamp": "20241015-132435", "method": "local", "text": "What are the top themes in this story?"}'
```

global:

```bash
curl -X POST localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{"kb": "raggo", "timestamp": "20241015-132435", "method": "global", "text": "Who is Scrooge, and what are his main relationships?"}'
```
