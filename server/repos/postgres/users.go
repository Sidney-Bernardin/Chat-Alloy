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
	PWSalt   []byte
}

const qSelectUserByID = `SELECT id, created_at, updated_at, username  FROM users WHERE id = $1`

func (repo *Repository) SelectUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := row[User](ctx, repo, qSelectUserByID, userID)
	return user, errors.WithStack(err)
}

const qSelectUserByPassword = `SELECT id, created_at, updated_at, username  FROM users WHERE pw_hash = $1`

func (repo *Repository) SelectUserByPassword(ctx context.Context, password string) (*User, error) {
	user, err := row[User](ctx, repo, qSelectUserByID, password)
	return user, errors.WithStack(err)
}

const qSelectUserSaltByUsername = `SELECT id, pw_salt FROM users WHERE username = $1`

func (repo *Repository) SelectUserSaltByUsername(ctx context.Context, username string) (uuid.UUID, []byte, error) {
	user, err := row[User](ctx, repo, qSelectUserByID, username)
	return user.ID, user.PWSalt, errors.WithStack(err)
}
