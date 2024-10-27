package pgx

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (s *SessionPGX) AddWithdraw(ctx context.Context, withdraw *domainmodels.Withdraw) (*domainmodels.Withdraw, error) {
	s.logger.Debug("Adding Withdraw", zap.Any("withdraw", withdraw))

	var id uint

	err := s.tx.QueryRow(ctx, `
		INSERT INTO public.withdrawals (order_number, user_uuid, sum, processed_at) VALUES 
			(@order_number, @user_uuid, @sum, @processed_at)
			RETURNING id
		`,
		pgx.NamedArgs{
			"order_number": withdraw.OrderNumber,
			"user_uuid":    withdraw.UserUUID,
			"sum":          withdraw.Sum.Float64(),
			"processed_at": withdraw.ProcessedAt,
		}).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) &&
			strings.Contains(pgErr.ConstraintName, "withdrawals_order_number_sum_user_uuid_unique") {
			return nil, domainmodels.ErrWithdrawExists
		}
		if pgErr.Code == pgerrcode.RaiseException &&
			strings.Contains(pgErr.Message, "Insufficient accrual for user") {
			return nil, domainmodels.ErrNotEnoughPointsToWithdraw
		}
	}
	if err != nil {
		s.logger.Error("Failed to insert into withdrawals",
			zap.Any("withdraw", withdraw),
			zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.AddWithdraw: %w", err))
	}
	withdraw.ID = id
	return withdraw, nil
}

func (s *SessionPGX) GetWithdrawalsByUserUUID(ctx context.Context, userUUID string) ([]*domainmodels.Withdraw, error) {
	s.logger.Debug("Getting user withdrawals", zap.Any("userUUID", userUUID))

	rows, err := s.tx.Query(ctx, `
	SELECT id, order_number, sum, processed_at 
	FROM public.withdrawals 
	WHERE user_uuid = @user_uuid 
	ORDER BY processed_at DESC
	`, pgx.NamedArgs{"user_uuid": userUUID})

	if err != nil {
		s.logger.Error("Failed to query withdrawals",
			zap.Any("userUUID", userUUID),
			zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetWithdrawalsByUserUUID: %w", err))
	}
	withdrawals := make([]*domainmodels.Withdraw, 0)
	for rows.Next() {
		var id uint
		var sum float64
		var orderNumber string
		var processedAt time.Time
		if err := rows.Scan(&id, &orderNumber, &sum, &processedAt); err != nil {
			s.logger.Error("Failed to scan row", zap.String("userUUID", userUUID), zap.Error(err))
			return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetWithdrawalsByUserUUID: %w", err))
		}
		withdrawal := &domainmodels.Withdraw{
			ID:          id,
			OrderNumber: orderNumber,
			Sum:         domainmodels.AccrualPointFromFloat64(sum),
			ProcessedAt: processedAt,
			UserUUID:    userUUID,
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	if err := rows.Err(); err != nil {
		s.logger.Error("Failed to query database", zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("error in SessionPGX.GetWithdrawalsByUserUUID: %w", err))
	}
	return withdrawals, nil
}
