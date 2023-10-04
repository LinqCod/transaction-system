package repository

import (
	"context"
	"database/sql"
	"github.com/linqcod/transaction-system/publisher_service/internal/model"
)

const (
	UpdateAccountFrozenBalanceQuery = `UPDATE accounts SET frozen_balance = $1 WHERE card_number=$2;`
	GetAccountQuery                 = `SELECT id, card_number, balance, frozen_balance FROM accounts WHERE card_number=$1`
	GetAccountsBalanceQuery         = `SELECT card_number, balance, frozen_balance FROM accounts;`
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

func (r AccountRepository) GetAccountsBalance() ([]*model.AccountBalance, error) {
	var accounts []*model.AccountBalance

	rows, err := r.db.QueryContext(r.ctx, GetAccountsBalanceQuery)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account model.AccountBalance
		if err := rows.Scan(
			&account.CardNumber,
			&account.ActualBalance,
			&account.FrozenBalance,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (r AccountRepository) UpdateAccountFrozenBalance(card string, newFrozenBalance float64) error {
	if err := r.db.QueryRowContext(r.ctx, UpdateAccountFrozenBalanceQuery, newFrozenBalance, card).Err(); err != nil {
		return err
	}

	return nil
}
