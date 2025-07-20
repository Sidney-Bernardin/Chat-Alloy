package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sidney-Bernardin/Chat-Alloy/internal"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CSRFToken string    `json:"csrf_token"`
}

func (repo *Repository) SetNewSession(ctx context.Context, userID uuid.UUID) (*Session, error) {

	var (
		session = &Session{
			ID:        uuid.New(),
			UserID:    userID,
			CSRFToken: internal.MustRandomString(10),
		}
		key = fmt.Sprintf("session:%s", session.ID)
	)

	err := repo.client.Watch(ctx, func(tx *redis.Tx) error {
		if err := repo.client.JSONSet(ctx, key, ".", session).Err(); err != nil {
			return errors.Wrap(err, "cannot set session")
		}

		err := repo.client.Expire(ctx, key, repo.cfg.SESSION_DURATION).Err()
		return errors.Wrap(err, "cannot set expire")
	}, key)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return session, nil
}

func (repo *Repository) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {

	key := fmt.Sprintf("session:%s", sessionID)
	sessionJSON, err := repo.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get session")
	}

	var session *Session
	if err := json.Unmarshal([]byte(sessionJSON), session); err != nil {
		return nil, errors.Wrap(err, "cannot decode session")
	}

	return session, nil
}
