package server

import "net/http"

type Dependencies struct {
	Chat       http.Handler
	Embeddings http.Handler
	Admin      http.Handler
}

func Routes(mux *http.ServeMux, d Dependencies) {
	mux.Handle("POST /v1/chat/completions", d.Chat)
	mux.Handle("POST /v1/embeddings", d.Embeddings)
	mux.Handle("GET /v1/admin/tenants", d.Admin)
}
