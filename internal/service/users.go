package service

import (
	"context"

	"github.com/Sidney-Bernardin/Chat-Alloy/internal"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/repos/postgres"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (svc *Service) Signup(ctx context.Context, username, password string) (sessionID uuid.UUID, err error) {

	pwSalt := internal.MustRandomString(16)
	pwHash, err := bcrypt.GenerateFromPassword([]byte(password+pwSalt), 12)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot hash password")
	}

	user := &postgres.User{
		ID:       uuid.New(),
		Username: username,
		PWHash:   pwHash,
		PWSalt:   pwSalt,
	}

	if err = svc.Postgres.InsertUser(ctx, user); err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot insert user")
	}

	session, err := svc.Redis.SetNewSession(ctx, user.ID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot set new session")
	}

	return session.ID, nil
}

func (svc *Service) Signin(ctx context.Context, username, password string) (sessionID uuid.UUID, err error) {

	userID, pwHash, pwSalt, err := svc.Postgres.SelectUserSaltByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return uuid.Nil, &DomainError{
				Type: DomainErrorTypeUserNotFound,
				Msg:  "User does not exist",
			}
		}

		return uuid.Nil, errors.Wrap(err, "cannot select user by username")
	}

	if err := bcrypt.CompareHashAndPassword(pwHash, []byte(password+pwSalt)); err != nil {
		return uuid.Nil, &DomainError{
			Type: DomainErrorTypeInvalidPassword,
		}
	}

	session, err := svc.Redis.SetNewSession(ctx, userID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot set new session")
	}

	return session.ID, nil
}
