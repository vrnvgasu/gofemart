package app

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/vrnvgasu/gofemart/internal/repository/postgres"
	"github.com/vrnvgasu/gofemart/pkg/jwt"
)

func (a *App) Login(ctx context.Context, login, password string) (string, error) {
	user, err := a.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("app.Login GetUserByLogin: %w", err)
	}
	if user == nil {
		return "", UnauthorizedError()
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", UnauthorizedError()
	}

	token, err := jwt.Generate(user.ID, a.cfg.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("app.Login Generate: %w", err)
	}

	return token, nil
}

func (a *App) Register(ctx context.Context, login, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("app.Register GenerateFromPassword: %w", err)
	}

	user, err := a.storage.CreateUser(ctx, login, string(hash))
	if err != nil {
		if errors.Is(err, postgres.ErrLoginConflict) {
			return "", ConflictError()
		}

		return "", fmt.Errorf("app.Register CreateUser: %w", err)
	}

	token, err := jwt.Generate(user.ID, a.cfg.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("app.Login Generate: %w", err)
	}

	return token, nil
}
