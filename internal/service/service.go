package service

import (
	"fmt"
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

type DomainError struct {
	Type  DomainErrorType `json:"type"`
	Msg   string          `json:"message,omitempty"`
	Attrs map[string]any  `json:"attrs,omitempty"`
}

type DomainErrorType string

var (
	DomainErrorTypeInvalidPassword DomainErrorType = "invalid-password"
	DomainErrorTypeUserNotFound    DomainErrorType = "user-not-found"
)

func (err *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", err.Type, err.Msg)
}
