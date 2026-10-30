# llmgate

OpenAI-compatible LLM gateway written in Go.

## Module

```
github.com/vviveksharma/llmgate
```

## Quick start

```bash
go run ./cmd/gateway
curl -s http://localhost:8080/healthz
curl -s -H "X-API-Key: dev-key" -H "Content-Type: application/json" \
  -d '{"model":"mock","messages":[{"role":"user","content":"hello"}]}' \
  http://localhost:8080/v1/chat/completions
```

## Layout

See `prompt.md` for the intended directory structure.

## Make targets

| Target       | Description              |
|-------------|--------------------------|
| `make build` | Build gateway and loadtest |
| `make run`   | Run the gateway          |
| `make test`  | Run tests                |
| `make loadtest` | Simple HTTP load test |
| `make docker-up` | Compose stack (Redis, Prometheus) |

## Configuration

Environment variables are documented in `internal/config/config.go`.
