package service

import (
	"log/slog"

	"github.com/Sidney-Bernardin/Chat-Alloy/server"
	"github.com/Sidney-Bernardin/Chat-Alloy/server/repos/postgres"
)

type Service struct {
	cfg *server.Config
	log *slog.Logger

	pg *postgres.Repository
}

func New(cfg *server.Config, log *slog.Logger, pg *postgres.Repository) *Service {
	return &Service{cfg, log, pg}
}
