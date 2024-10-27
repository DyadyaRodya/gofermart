package pgx

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (s *SessionPGX) GetTransactionSumInfoByUserUUID(ctx context.Context, userUUID string) (*interactorsdto.TransactionSumInfo, error) {
	s.logger.Debug("Getting user transactions sum info", zap.Any("userUUID", userUUID))

	var totalAccrual, totalWithdrawn float64
	err := s.tx.QueryRow(ctx, `WITH total_accruals as (SELECT
		u.uuid as user_uuid,
		COALESCE(SUM(o.accrual), 0) as total_accrual
	FROM public.users AS u
		LEFT JOIN public.orders AS o ON u.uuid = o.user_uuid
	WHERE u.uuid = @user_uuid
	GROUP BY u.uuid),
	total_withdrawals as (SELECT
		u.uuid as user_uuid,
		COALESCE(SUM(w.sum), 0) as total_withdawn
	FROM public.users AS u
		LEFT JOIN public.withdrawals AS w on u.uuid = w.user_uuid
	WHERE u.uuid = @user_uuid
	GROUP BY u.uuid)
	SELECT total_accrual, total_withdawn
	FROM total_accruals AS ta
    	JOIN total_withdrawals AS tw ON ta.user_uuid=tw.user_uuid;
	`, pgx.NamedArgs{"user_uuid": userUUID}).Scan(&totalAccrual, &totalWithdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainmodels.ErrUserNotFound
		}
		s.logger.Error("Failed to query transactions sum info",
			zap.Any("userUUID", userUUID),
			zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetTransactionSumInfoByUserUUID: %w", err))
	}

	return &interactorsdto.TransactionSumInfo{
		OrderAccruals: domainmodels.AccrualPointFromFloat64(totalAccrual),
		Withdrawn:     domainmodels.AccrualPointFromFloat64(totalWithdrawn),
	}, nil
}
