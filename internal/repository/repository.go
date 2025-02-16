//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
package repository

import (
	"Avito/internal/db"
	"Avito/internal/model"
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type Repo struct {
	db db.DBops
}

type Repository interface {
	BeginTransaction(ctx context.Context, options *pgx.TxOptions) (pgx.Tx, error)
	GetUser(ctx context.Context, tx pgx.Tx, login string) (*User, error)
	AddUser(ctx context.Context, tx pgx.Tx, login string, password string) error
	UpdateUserBalance(ctx context.Context, tx pgx.Tx, login string, amount int32) error
	GetUserItems(ctx context.Context, tx pgx.Tx, login string) ([]model.Item, error)
	GetReceiveCoinEvents(ctx context.Context, tx pgx.Tx, login string) ([]model.ReceiveCoinEvent, error)
	GetSendCoinEvents(ctx context.Context, tx pgx.Tx, login string) ([]model.SendCoinEvent, error)
	AddSendCoinEvent(ctx context.Context, tx pgx.Tx, fromUser string, toUser string, amount int32) error
	UpdateItemPurchaseCount(ctx context.Context, tx pgx.Tx, login string, item string) error
	CommitTx(ctx context.Context, tx pgx.Tx)
	RollbackTx(ctx context.Context, tx pgx.Tx)
}

func NewRepository(database db.DBops) *Repo {
	return &Repo{db: database}
}

func (r *Repo) BeginTransaction(ctx context.Context, options *pgx.TxOptions) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, options)
}

func (r *Repo) RollbackTx(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil {
		log.Println("failed to rollback tx wih err:", err)
	}
}

func (r *Repo) CommitTx(ctx context.Context, tx pgx.Tx) {
	if err := tx.Commit(ctx); err != nil {
		log.Println("failed to commit tx wih err:", err)
	}
}

func (r *Repo) GetUser(ctx context.Context, tx pgx.Tx, login string) (*User, error) {
	var user User
	err := tx.QueryRow(ctx, "SELECT login, password, coins FROM users WHERE login=$1 FOR UPDATE", login).Scan(&user.Login, &user.Password, &user.Coins)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repo) AddUser(ctx context.Context, tx pgx.Tx, login string, password string) error {
	_, err := tx.Exec(ctx, "INSERT INTO users (login, password, coins) VALUES ($1, $2, $3)", login, password, 1000)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateUserBalance(ctx context.Context, tx pgx.Tx, login string, amount int32) error {
	_, err := tx.Exec(ctx, "UPDATE users SET coins = coins + $1 WHERE login = $2", amount, login)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateItemPurchaseCount(ctx context.Context, tx pgx.Tx, login string, item string) error {
	_, err := tx.Exec(ctx, `INSERT INTO purchases (login, type, quantity) VALUES ($1, $2, 1) 
		ON CONFLICT (login, type) DO UPDATE SET quantity = purchases.quantity + 1`, login, item)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) AddSendCoinEvent(ctx context.Context, tx pgx.Tx, fromUser string, toUser string, amount int32) error {
	_, err := tx.Exec(ctx, "INSERT INTO send_coin_events (to_user, from_user, amount) VALUES ($1, $2, $3)", toUser, fromUser, amount)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetSendCoinEvents(ctx context.Context, tx pgx.Tx, login string) ([]model.SendCoinEvent, error) {
	rows, err := tx.Query(ctx, `SELECT to_user, SUM(amount) FROM send_coin_events WHERE from_user=$1 GROUP BY to_user`, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []model.SendCoinEvent{}
	for rows.Next() {
		var event model.SendCoinEvent
		err := rows.Scan(&event.ToUser, &event.Amount)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return events, nil
}

func (r *Repo) GetReceiveCoinEvents(ctx context.Context, tx pgx.Tx, login string) ([]model.ReceiveCoinEvent, error) {
	rows, err := tx.Query(ctx, `SELECT from_user, SUM(amount) FROM send_coin_events WHERE to_user=$1 GROUP BY from_user`, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []model.ReceiveCoinEvent{}
	for rows.Next() {
		var event model.ReceiveCoinEvent
		err := rows.Scan(&event.FromUser, &event.Amount)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return events, nil
}

func (r *Repo) GetUserItems(ctx context.Context, tx pgx.Tx, login string) ([]model.Item, error) {
	rows, err := tx.Query(ctx, `SELECT type, quantity FROM purchases WHERE login=$1`, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Item{}
	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.Type, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return items, nil
}
