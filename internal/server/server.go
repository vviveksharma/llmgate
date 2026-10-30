package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/vviveksharma/llmgate/internal/auth"
	"github.com/vviveksharma/llmgate/internal/config"
	"github.com/vviveksharma/llmgate/internal/handler"
	"github.com/vviveksharma/llmgate/internal/middleware"
	"github.com/vviveksharma/llmgate/internal/observability"
	"github.com/vviveksharma/llmgate/internal/provider"
	"github.com/vviveksharma/llmgate/internal/router"
	"github.com/vviveksharma/llmgate/internal/storage"
	"github.com/vviveksharma/llmgate/internal/tenant"
	"github.com/vviveksharma/llmgate/internal/usage"
)

type Server struct {
	cfg    *config.Config
	log    *slog.Logger
	http   *http.Server
	pg     *storage.Postgres
	redis  *storage.Redis
}

func New(cfg *config.Config, log *slog.Logger) (*Server, error) {
	pg, err := storage.NewPostgres(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	rds := storage.NewRedis(cfg.RedisAddr)

	reg := provider.NewRegistry(cfg.DefaultProvider)
	reg.Register(provider.NewMock())
	reg.Register(provider.NewOllama(""))

	modelRouter := router.NewModelRouter(reg)
	usageRepo := usage.NewMemoryRepository()
	recorder := usage.NewRecorder(usageRepo)

	tenantSvc := tenant.NewService(tenant.NewMemoryRepository(nil))

	apiKeys := auth.NewAPIKeyValidator(map[string]string{"dev-key": "tenant-1"})
	jwtVal := auth.NewJWTValidator(cfg.JWTSecret)

	root := http.NewServeMux()
	root.Handle("GET /metrics", observability.MetricsHandler())
	root.Handle("GET /healthz", &handler.Health{Ready: func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := pg.Ping(ctx); err != nil {
			return err
		}
		return rds.Ping(ctx)
	}})

	api := http.NewServeMux()
	Routes(api, Dependencies{
		Chat:       &handler.Chat{Router: modelRouter, Recorder: recorder},
		Embeddings: &handler.Embeddings{Router: modelRouter},
		Admin:      &handler.Admin{Tenants: tenantSvc, JWT: jwtVal},
	})

	var protected http.Handler = api
	protected = middleware.NewRateLimit(cfg.RateLimitRPS, cfg.RateLimitBurst).Middleware(protected)
	protected = middleware.Auth{APIKeys: apiKeys, JWT: jwtVal}.Middleware(protected)
	root.Handle("/", protected)

	var h http.Handler = root
	h = middleware.Metrics(h)
	h = middleware.RequestID(h)
	h = middleware.Logging(log)(h)
	h = middleware.Recover(log)(h)

	s := &Server{
		cfg:   cfg,
		log:   log,
		pg:    pg,
		redis: rds,
		http: &http.Server{
			Addr:    cfg.HTTPAddr,
			Handler: h,
		},
	}

	return s, nil
}

func (s *Server) Start() error {
	s.log.Info("starting gateway", "addr", s.cfg.HTTPAddr)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	_ = s.pg.Close()
	_ = s.redis.Close()
	return s.http.Shutdown(ctx)
}
