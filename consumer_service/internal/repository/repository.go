package repository

import (
	"context"
	"database/sql"
	"github.com/linqcod/transaction-system/consumer_service/internal/model"
)

const (
	UpdateAccountFrozenBalanceQuery = `UPDATE accounts SET frozen_balance = $1 WHERE card_number=$2;`
	UpdateAccountBalanceQuery       = `UPDATE accounts SET balance = $1 WHERE card_number=$2;`
	GetAccountQuery                 = `SELECT id, card_number, balance, frozen_balance FROM accounts WHERE card_number=$1`
)

type AccountRepository struct {
	ctx context.Context
	db  *sql.DB
}

func NewAccountRepository(ctx context.Context, db *sql.DB) *AccountRepository {
	return &AccountRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r AccountRepository) GetAccount(card string) (*model.Account, error) {
	var account model.Account

	if err := r.db.QueryRowContext(r.ctx, GetAccountQuery, card).Scan(
		&account.Id,
		&account.CardNumber,
		&account.Balance,
		&account.FrozenBalance,
	); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r AccountRepository) UpdateAccountBalance(card string, newBalance float64) error {
	if err := r.db.QueryRowContext(r.ctx, UpdateAccountBalanceQuery, newBalance, card).Err(); err != nil {
		return err
	}

	return nil
}

func (r AccountRepository) UpdateAccountFrozenBalance(card string, newFrozenBalance float64) error {
	if err := r.db.QueryRowContext(r.ctx, UpdateAccountFrozenBalanceQuery, newFrozenBalance, card).Err(); err != nil {
		return err
	}

	return nil
}
