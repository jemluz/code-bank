package repository

import (
	"database/sql"
	"errors"

	"github.com/codeedu/codebank/domain"
)

type TransactionRepositoryDB struct {
	db *sql.DB
}

func (t *TransactionRepositoryDB) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	stmt, err := t.db.Prepare("select id, balance, balance_limit from credit_cards where number=$1")

	if err != nil {
		return c, err
	}

	if err = stmt.QueryRow(creditCard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card does not exists")
	}
	return c, nil
}

func NewTransactionRepositoryDB(db *sql.DB) *TransactionRepositoryDB {
	return &TransactionRepositoryDB{db: db}
}

func (t *TransactionRepositoryDB) SaveTransaction(transaction domain.Transaction, creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(
		`insert into transactions(id, credit_card_id, amount, status, descriptio, store, created_at) values($1, $2, $3, $4, $5, $6, $7)`,
	)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transaction.ID,
		transaction.CreditCardId,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)
	if err != nil {
		return err
	}

	if transaction.Status == "approved" {
		err = t.updateBalance(creditCard)
		if err != nil {
			return err
		}
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepositoryDB) updateBalance(creditCard domain.CreditCard) error {
	_, err := t.db.Exec(
		"update credit_cards set balance = $1 where id = $2",
		creditCard.Balance,
		creditCard.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDB) CreateCreditCard(creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(
		`insert into transactions(id, name, number, expiration_month, expiration_yeah, CVV, balance, balance_limit, created_at) values($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		creditCard.ID,
		creditCard.Name,
		creditCard.Number,
		creditCard.ExpirationMonth,
		creditCard.ExpirationYear,
		creditCard.CVV,
		creditCard.Balance,
		creditCard.Limit,
	)

	if err != nil {
		return err
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}
