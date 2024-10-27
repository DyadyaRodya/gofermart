package pgx

import (
	"context"
	"database/sql"
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

func (s *SessionPGX) AddOrder(ctx context.Context, order *domainmodels.Order) (*domainmodels.Order, error) {
	s.logger.Debug("Adding Order", zap.Any("order", order))

	var id uint
	accrual := sql.NullFloat64{
		Float64: order.Accrual.Float64(),
		Valid:   order.Accrual.Float64() > 0,
	}

	err := s.tx.QueryRow(ctx, `
		INSERT INTO public.orders (number, user_uuid, status, uploaded_at, accrual) VALUES 
			(@number, @user_uuid, @status, @uploaded_at, @accrual)
			RETURNING id
		`,
		pgx.NamedArgs{
			"number":      order.Number,
			"user_uuid":   order.UserUUID,
			"status":      order.Status,
			"uploaded_at": order.UploadedAt,
			"accrual":     accrual,
		}).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) &&
		pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) &&
		strings.Contains(pgErr.ConstraintName, "orders_number_key") {
		return nil, domainmodels.ErrOrderExists
	}

	if err != nil {
		s.logger.Error("Failed to insert into orders",
			zap.Any("order", order),
			zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.AddOrder: %w", err))
	}
	order.ID = id
	return order, nil
}

func (s *SessionPGX) UpdateOrder(ctx context.Context, order *domainmodels.Order) error {
	s.logger.Debug("Updating Order", zap.Any("order", order))

	accrual := sql.NullFloat64{
		Float64: order.Accrual.Float64(),
		Valid:   order.Accrual.Float64() > 0,
	}

	_, err := s.tx.Exec(ctx, `UPDATE public.orders SET status = @status, accrual = @accrual WHERE id = @id`,
		pgx.NamedArgs{
			"id":      order.ID,
			"status":  order.Status,
			"accrual": accrual,
		})

	if err != nil {
		s.logger.Error("Failed to update order",
			zap.Any("order", order),
			zap.Error(err))
		return errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.UpdateOrder: %w", err))
	}
	return nil
}

func (s *SessionPGX) GetOrdersByUserUUID(ctx context.Context, userUUID string) ([]*domainmodels.Order, error) {
	s.logger.Debug("Getting user orders", zap.Any("userUUID", userUUID))

	rows, err := s.tx.Query(ctx, `
	SELECT id, number, status, uploaded_at, accrual 
	FROM public.orders 
	WHERE user_uuid = @user_uuid 
	ORDER BY uploaded_at DESC
	`, pgx.NamedArgs{"user_uuid": userUUID})

	if err != nil {
		s.logger.Error("Failed to query orders",
			zap.Any("userUUID", userUUID),
			zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetOrdersByUserUUID: %w", err))
	}
	orders := make([]*domainmodels.Order, 0)
	for rows.Next() {
		var id uint
		var accrualSQL sql.NullFloat64
		var number string
		var status domainmodels.OrderStatus
		var uploadedAt time.Time
		if err := rows.Scan(&id, &number, &status, &uploadedAt, &accrualSQL); err != nil {
			s.logger.Error("Failed to scan row", zap.String("userUUID", userUUID), zap.Error(err))
			return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetOrdersByUserUUID: %w", err))
		}
		var accrual domainmodels.AccrualPoint
		if accrualSQL.Valid {
			accrual = domainmodels.AccrualPointFromFloat64(accrualSQL.Float64)
		} else {
			accrual = 0
		}
		order := &domainmodels.Order{
			ID:         id,
			Number:     number,
			Status:     status,
			UploadedAt: uploadedAt,
			Accrual:    accrual,
			UserUUID:   userUUID,
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		s.logger.Error("Failed to query database", zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("error in SessionPGX.GetOrdersByUserUUID: %w", err))
	}
	return orders, nil
}

func (s *SessionPGX) GetOrderByNumberAndUserUUID(
	ctx context.Context, number, userUUID string,
) (*domainmodels.Order, error) {
	s.logger.Debug("Getting user orders", zap.String("number", number), zap.String("userUUID", userUUID))

	var id uint
	var accrualSQL sql.NullFloat64
	var status domainmodels.OrderStatus
	var uploadedAt time.Time
	err := s.tx.QueryRow(ctx, `
	SELECT id, status, uploaded_at, accrual 
	FROM public.orders 
	WHERE user_uuid = @user_uuid AND number = @number
	ORDER BY uploaded_at DESC
	`, pgx.NamedArgs{"number": number, "user_uuid": userUUID}).
		Scan(&id, &status, &uploadedAt, &accrualSQL)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		s.logger.Error("Failed to query order",
			zap.String("number", number),
			zap.String("userUUID", userUUID),
			zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetOrderByNumberAndUserUUID: %w", err))
	}
	var accrual domainmodels.AccrualPoint
	if accrualSQL.Valid {
		accrual = domainmodels.AccrualPointFromFloat64(accrualSQL.Float64)
	} else {
		accrual = 0
	}
	order := &domainmodels.Order{
		ID:         id,
		Number:     number,
		Status:     status,
		UploadedAt: uploadedAt,
		Accrual:    accrual,
		UserUUID:   userUUID,
	}
	return order, nil
}

func (s *SessionPGX) GetUnprocessedOrders(ctx context.Context) ([]*domainmodels.Order, error) {
	s.logger.Debug("Getting unprocessed orders")

	rows, err := s.tx.Query(ctx, `
	SELECT id, number, status, uploaded_at, accrual, user_uuid
	FROM public.orders 
	WHERE status in ('NEW', 'PROCESSING') 
	ORDER BY uploaded_at
	`) // order ASC to process older first

	if err != nil {
		s.logger.Error("Failed to query unprocessed orders",
			zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetUnprocessedOrders: %w", err))
	}
	orders := make([]*domainmodels.Order, 0)
	for rows.Next() {
		var id uint
		var accrualSQL sql.NullFloat64
		var number, userUUID string
		var status domainmodels.OrderStatus
		var uploadedAt time.Time
		if err := rows.Scan(&id, &number, &status, &uploadedAt, &accrualSQL, &userUUID); err != nil {
			s.logger.Error("Failed to scan row", zap.Error(err))
			return nil, errors.Join(ErrDBAPI, fmt.Errorf("SessionPGX.GetUnprocessedOrders: %w", err))
		}
		var accrual domainmodels.AccrualPoint
		if accrualSQL.Valid {
			accrual = domainmodels.AccrualPointFromFloat64(accrualSQL.Float64)
		} else {
			accrual = 0
		}
		order := &domainmodels.Order{
			ID:         id,
			Number:     number,
			Status:     status,
			UploadedAt: uploadedAt,
			Accrual:    accrual,
			UserUUID:   userUUID,
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		s.logger.Error("Failed to query database", zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("error in SessionPGX.GetUnprocessedOrders: %w", err))
	}
	return orders, nil
}
