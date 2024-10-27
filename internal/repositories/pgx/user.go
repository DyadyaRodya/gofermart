package pgx

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"time"
)

func (s *SessionPGX) GetUserByLogin(ctx context.Context, login string) (*domainmodels.UserInfo, error) {
	s.logger.Debug("GetUserByLogin", zap.String("login", login))

	var uuid, passwordHash, passwordSalt string
	var createdAt time.Time
	err := s.tx.QueryRow(ctx, `SELECT uuid, created_at, password_hash, password_salt FROM public.users 
                                                      WHERE login = @login AND deleted_at IS NULL`,
		pgx.NamedArgs{"login": login},
	).Scan(&uuid, &createdAt, &passwordHash, &passwordSalt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainmodels.ErrUserNotFound
	}
	if err != nil {
		s.logger.Error("Failed to get user by login", zap.String("login", login), zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetUserByLogin: %w", err))
	}

	return &domainmodels.UserInfo{
		UUID:         uuid,
		Login:        login,
		CreatedAt:    createdAt,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}, nil
}

func (s *SessionPGX) AddUser(ctx context.Context, user *domainmodels.UserInfo) error {
	s.logger.Debug("Adding User", zap.Any("user", user))

	ct, err := s.tx.Exec(ctx, `
		INSERT INTO public.users (uuid, login, created_at, password_hash, password_salt) VALUES 
			(@uuid, @login, @created_at, @password_hash, @password_salt)
		`,
		pgx.NamedArgs{
			"uuid":          user.UUID,
			"login":         user.Login,
			"created_at":    user.CreatedAt,
			"password_hash": user.PasswordHash,
			"password_salt": user.PasswordSalt,
		})

	if err != nil {
		s.logger.Error("Failed to insert into users",
			zap.Any("user", user),
			zap.Error(err))
		return errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.AddUser: %w", err))
	}
	if !ct.Insert() {
		s.logger.Error("Failed to insert into users",
			zap.Any("user", user),
			zap.Any("commandTag", ct))
		return errors.Join(ErrDBAPI, fmt.Errorf("error in SessionPGX.AddUser: not inserted user: %v", user))
	}
	return nil
}
