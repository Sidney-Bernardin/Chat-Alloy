package service

import (
	"log/slog"

	"github.com/Sidney-Bernardin/Chat-Alloy/internal"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/repos/postgres"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/repos/redis"
)

type Service struct {
	Config *internal.Config
	Logger *slog.Logger

	Postgres *postgres.Repository
	Redis    *redis.Repository
}
