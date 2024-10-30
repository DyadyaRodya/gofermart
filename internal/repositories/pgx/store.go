package pgx

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var ErrDBAPI = errors.Join(domainmodels.ErrInternalServer, errors.New("database api error"))

type SessionPGX struct {
	tx     pgx.Tx
	logger *zap.Logger
}

func (s *SessionPGX) Commit(ctx context.Context) error {
	err := s.tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Failed to Commit", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}
	return nil
}

func (s *SessionPGX) Close(ctx context.Context) error {
	err := s.tx.Rollback(ctx)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		s.logger.Error("Failed to Rollback", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}
	return nil
}

type StorePGX struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewStorePGX(pool *pgxpool.Pool, logger *zap.Logger) *StorePGX {
	return &StorePGX{pool: pool, logger: logger}
}

func (s *StorePGX) InitSchema(ctx context.Context) error {
	s.logger.Info("Initializing schema")
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}
	defer tx.Rollback(ctx)

	s.logger.Info("Creating table `users`")
	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS public.users (
        uuid UUID NOT NULL PRIMARY KEY, 
        login VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
		password_hash VARCHAR(255),
    	password_salt VARCHAR(255)
    )`)
	if err != nil {
		s.logger.Error("Failed to create table `users`", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}

	s.logger.Info("Creating type `order_status`")
	_, err = tx.Exec(ctx, `DO $$ BEGIN
    CREATE TYPE public.order_status AS ENUM (
        'NEW', 'PROCESSING', 'INVALID', 'PROCESSED'
	);
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;
	`)
	if err != nil {
		s.logger.Error("Failed to create type `order_status`", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}

	s.logger.Info("Creating table `orders`")
	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS public.orders (
    	id SERIAL PRIMARY KEY,
    	number VARCHAR(255) UNIQUE NOT NULL,
        user_uuid UUID REFERENCES public.users(uuid) ON DELETE CASCADE,
        status public.order_status,
        uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        accrual DECIMAL(12,2) NULL DEFAULT NULL
	)`)
	if err != nil {
		s.logger.Error("Failed to create table `orders`", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}

	s.logger.Info("Creating table `withdrawals`")
	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS public.withdrawals (
    	id SERIAL PRIMARY KEY,
    	order_number VARCHAR(255) NOT NULL,
        user_uuid UUID REFERENCES public.users(uuid) ON DELETE CASCADE,
        sum DECIMAL(12,2) NULL,
        processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        CONSTRAINT withdrawals_order_number_sum_unique UNIQUE (order_number, sum)
	)`)
	if err != nil {
		s.logger.Error("Failed to create table `withdrawals`", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}

	s.logger.Info("Creating function `check_withdrawal_sum`")
	_, err = tx.Exec(ctx, `
	CREATE OR REPLACE FUNCTION public.check_withdrawal_sum()
		RETURNS TRIGGER AS $$
	DECLARE
		total_accrual INTEGER;
		total_withdrawal INTEGER;
	BEGIN
		-- Calculate the total accrual for the user
		SELECT COALESCE(SUM(COALESCE(accrual, 0)), 0)
		INTO total_accrual
		FROM public.orders
		WHERE user_uuid = NEW.user_uuid;
	
		-- Calculate the total withdrawal for the user before the new insert
		SELECT COALESCE(SUM(sum), 0)
		INTO total_withdrawal
		FROM public.withdrawals
		WHERE user_uuid = NEW.user_uuid;
	
		-- Check if the total accrual is greater than or equal to total withdrawal + new sum
		IF total_accrual < (total_withdrawal + NEW.sum) THEN
			RAISE EXCEPTION 'Insufficient accrual for user %: Total accrual % is less than required %',
				NEW.user_uuid, total_accrual, (total_withdrawal + NEW.sum);
		END IF;
	
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		s.logger.Error("Failed to create function `check_withdrawal_sum`", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}

	s.logger.Info("Creating trigger `check_withdrawal_sum`")
	_, err = tx.Exec(ctx, `CREATE OR REPLACE TRIGGER validate_withdrawal
		BEFORE INSERT ON withdrawals
		FOR EACH ROW EXECUTE FUNCTION check_withdrawal_sum();
	`)
	if err != nil {
		s.logger.Error("Failed to create trigger `check_withdrawal_sum`", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return errors.Join(ErrDBAPI, err)
	}
	s.logger.Info("Initializing schema done")
	return nil
}

func (s *StorePGX) NewSession(ctx context.Context) (interfaces.RepoSession, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("error in StorePGX.NewSession: %w", err))
	}
	txPGX := &SessionPGX{tx: tx, logger: s.logger}
	return txPGX, nil
}

func (s *StorePGX) NewSerializableSession(ctx context.Context) (interfaces.RepoSession, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("error in StorePGX.NewSerializableSession: %w", err))
	}
	_, err = tx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		s.logger.Error("Failed to set transaction mode SERIALIZABLE", zap.Error(err))
		return nil, errors.Join(ErrDBAPI, fmt.Errorf("error in StorePGX.NewSerializableSession: %w", err))
	}
	txPGX := &SessionPGX{tx: tx, logger: s.logger}
	return txPGX, nil
}
