# Design decisions

- **Go module**: `github.com/vviveksharma/llmgate`
- **stdlib HTTP** with Go 1.22+ method patterns on `ServeMux` for minimal dependencies in the skeleton
- **Mock provider** as default for local development without external LLM services
- **API key auth** (`X-API-Key`) and optional JWT for multi-tenant scenarios
- **In-memory** tenant and usage stores in the scaffold; Postgres migrations live under `migrations/` for future persistence
