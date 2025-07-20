package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type User struct {
	ID uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time

	Username string
	PWHash   []byte
	PWSalt   string
}

const qSelectUserByID = `SELECT id, created_at, updated_at, username  FROM users WHERE id = $1`

func (repo *Repository) SelectUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := row[User](ctx, repo, qSelectUserByID, userID)
	return user, errors.WithStack(err)
}

const qSelectUserSaltByUsername = `SELECT id, pw_hash, pw_salt FROM users WHERE username = $1`

func (repo *Repository) SelectUserSaltByUsername(ctx context.Context, username string) (userID uuid.UUID, pwHash []byte, pwSalt string, err error) {
	user, err := row[User](ctx, repo, qSelectUserByID, username)
	return user.ID, user.PWHash, user.PWSalt, errors.WithStack(err)
}
