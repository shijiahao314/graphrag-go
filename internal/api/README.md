# Api

## kb

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