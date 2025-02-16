//go:generate mockgen -source ./service.go -destination=./mocks/service.go -package=mock_service
package service

import (
	"Avito/internal/errors"
	"Avito/internal/model"
	"Avito/internal/repository"
	"Avito/internal/store"
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type Svc struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Svc {
	return &Svc{
		repo: repo,
	}
}

type Service interface {
	GetUser(ctx context.Context, login string, password string) error
	MakePurchase(ctx context.Context, login string, item string) error
	GetUserInfo(ctx context.Context, login string) (*model.InfoResponse, error)
	SendCoin(ctx context.Context, fromUser string, toUser string, amount int32) error
}

func (s *Svc) addUser(ctx context.Context, login string, password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return err
	}

	user, err := s.repo.GetUser(ctx, tx, login)
	if err != nil {
		log.Println("failed to get user with err:", err)
		s.repo.RollbackTx(ctx, tx)
		return err
	}

	if user == nil {
		err = s.repo.AddUser(ctx, tx, login, string(hashPassword))
		if err != nil {
			log.Println("failed to add user with err:", err)
			s.repo.RollbackTx(ctx, tx)
			return err
		}
		s.repo.CommitTx(ctx, tx)
		return nil
	}
	s.repo.CommitTx(ctx, tx)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.ErrAuthFailed
	}

	return nil
}

func (s *Svc) GetUser(ctx context.Context, login string, password string) error {
	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return err
	}

	user, err := s.repo.GetUser(ctx, tx, login)
	if err != nil {
		log.Println("failed to get user with err:", err)
		s.repo.RollbackTx(ctx, tx)
		return err
	}
	s.repo.CommitTx(ctx, tx)

	if user == nil {
		return s.addUser(ctx, login, password)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.ErrAuthFailed
	}

	return nil
}

func (s *Svc) MakePurchase(ctx context.Context, login string, item string) error {
	price, ok := store.Store[item]
	if !ok {
		return errors.ErrItemNotFound
	}

	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return err
	}

	defer func() {
		if err != nil {
			s.repo.RollbackTx(ctx, tx)
		} else {
			s.repo.CommitTx(ctx, tx)
		}
	}()

	user, err := s.repo.GetUser(ctx, tx, login)
	if err != nil {
		log.Println("failed to get user with err:", err)
		return err
	}

	if user == nil {
		return errors.ErrUserDoesNotExist
	}

	if user.Coins < price {
		return errors.ErrNotEnoughMoney
	}

	err = s.repo.UpdateUserBalance(ctx, tx, login, -price)
	if err != nil {
		log.Println("failed to update user balance with err:", err)
		return err
	}

	err = s.repo.UpdateItemPurchaseCount(ctx, tx, login, item)
	if err != nil {
		log.Println("failed to update item purchase count with err:", err)
		return err
	}

	return nil
}

func (s *Svc) SendCoin(ctx context.Context, fromUser string, toUser string, amount int32) error {
	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return err
	}

	defer func() {
		if err != nil {
			s.repo.RollbackTx(ctx, tx)
		} else {
			s.repo.CommitTx(ctx, tx)
		}
	}()

	receiver, err := s.repo.GetUser(ctx, tx, toUser)
	if err != nil {
		log.Println("failed to get user receiver with err:", err)
		return err
	}

	if receiver == nil {
		return errors.ErrReceiverNotFound
	}

	sender, err := s.repo.GetUser(ctx, tx, fromUser)
	if err != nil {
		log.Println("failed to get user sender with err:", err)
		return err
	}

	if sender == nil {
		return errors.ErrUserDoesNotExist
	}

	if sender.Coins < amount {
		return errors.ErrNotEnoughMoney
	}

	err = s.repo.UpdateUserBalance(ctx, tx, fromUser, -amount)
	if err != nil {
		log.Println("failed to update user sender balance with err:", err)
		return err
	}

	err = s.repo.UpdateUserBalance(ctx, tx, toUser, amount)
	if err != nil {
		log.Println("failed to update user receiver balance with err:", err)
		return err
	}

	err = s.repo.AddSendCoinEvent(ctx, tx, fromUser, toUser, amount)
	if err != nil {
		log.Println("failed to add send coin event with err:", err)
		return err
	}

	return nil
}

func (s *Svc) GetUserInfo(ctx context.Context, login string) (*model.InfoResponse, error) {
	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			s.repo.RollbackTx(ctx, tx)
		} else {
			s.repo.CommitTx(ctx, tx)
		}
	}()

	user, err := s.repo.GetUser(ctx, tx, login)
	if err != nil {
		log.Println("failed to get user receiver with err:", err)
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserDoesNotExist
	}

	sent, err := s.repo.GetSendCoinEvents(ctx, tx, login)
	if err != nil {
		log.Println("failed to get send coin events with err:", err)
		return nil, err
	}

	received, err := s.repo.GetReceiveCoinEvents(ctx, tx, login)
	if err != nil {
		log.Println("failed to get received coin events with err:", err)
		return nil, err
	}

	items, err := s.repo.GetUserItems(ctx, tx, login)
	if err != nil {
		log.Println("failed to get user items with err:", err)
		return nil, err
	}

	return &model.InfoResponse{
		Coins: user.Coins,
		CoinHistory: model.CoinHistory{
			Received: received,
			Sent:     sent,
		},
		Inventory: items,
	}, nil
}
