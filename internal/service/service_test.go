package service

import (
	"context"
	"errors"
	"testing"

	customErrors "Avito/internal/errors"
	"Avito/internal/model"
	"Avito/internal/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func Test_GetUser(t *testing.T) {
	t.Parallel()
	var (
		ctx               = context.Background()
		correctPassword   = "password"
		incorrectPassword = "qwerty"
		login             = "user"
		hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
		user              = &repository.User{
			Password: string(hashedPassword),
		}
		userReqEarlier = &repository.User{
			Password: "another_password",
		}
	)

	t.Run("auth user", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)

		err := s.svc.GetUser(ctx, login, correctPassword)

		require.NoError(t, err)
	})

	t.Run("reg user", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return().AnyTimes()
		s.mockRepo.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := s.svc.GetUser(ctx, login, correctPassword)

		require.NoError(t, err)
	})

	t.Run("auth with wrong password", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.GetUser(ctx, login, incorrectPassword)

		require.EqualError(t, err, customErrors.ErrAuthFailed.Error())
	})

	t.Run("reg when before another user has reg", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return().AnyTimes()
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(userReqEarlier, nil)

		err := s.svc.GetUser(ctx, login, correctPassword)

		require.EqualError(t, err, customErrors.ErrAuthFailed.Error())
	})
}

func Test_MakePurchase(t *testing.T) {
	t.Parallel()

	var (
		ctx          = context.Background()
		login        = "user"
		item         = "t-shirt"
		notFoundItem = "item"
		user         = &repository.User{
			Coins: 100,
		}
		userHasNoMoney = &repository.User{
			Coins: 0,
		}
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)
		s.mockRepo.EXPECT().UpdateUserBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.mockRepo.EXPECT().UpdateItemPurchaseCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.MakePurchase(ctx, login, item)

		require.NoError(t, err)
	})

	t.Run("item not found", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		err := s.svc.MakePurchase(ctx, login, notFoundItem)

		require.EqualError(t, err, customErrors.ErrItemNotFound.Error())
	})

	t.Run("not enough money", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(userHasNoMoney, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.MakePurchase(ctx, login, item)

		require.EqualError(t, err, customErrors.ErrNotEnoughMoney.Error())
	})
}

func Test_SendCoin(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		senderLogin   = "sender"
		receiverLogin = "receiver"
		sender        = &repository.User{
			Coins: 100,
		}
		receiver    = &repository.User{}
		errorUpdate = errors.New("failed to update")
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(receiver, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(sender, nil)
		s.mockRepo.EXPECT().UpdateUserBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		s.mockRepo.EXPECT().AddSendCoinEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.SendCoin(ctx, senderLogin, receiverLogin, 10)

		require.NoError(t, err)
	})

	t.Run("receiver not found", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.SendCoin(ctx, senderLogin, receiverLogin, 10)

		require.EqualError(t, err, customErrors.ErrReceiverNotFound.Error())
	})

	t.Run("not enough money", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(receiver, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(sender, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.SendCoin(ctx, senderLogin, receiverLogin, 10000)

		require.EqualError(t, err, customErrors.ErrNotEnoughMoney.Error())
	})

	t.Run("failed update", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(receiver, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(sender, nil)
		s.mockRepo.EXPECT().UpdateUserBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.mockRepo.EXPECT().UpdateUserBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(errorUpdate.Error()))
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Times(1)

		err := s.svc.SendCoin(ctx, senderLogin, receiverLogin, 10)

		require.EqualError(t, err, errorUpdate.Error())
	})
}

func Test_GetUserInfo(t *testing.T) {
	t.Parallel()

	var (
		ctx      = context.Background()
		login    = "user"
		userInfo = &model.InfoResponse{
			Coins:     0,
			Inventory: []model.Item{},
			CoinHistory: model.CoinHistory{
				Sent:     []model.SendCoinEvent{},
				Received: []model.ReceiveCoinEvent{},
			},
		}
		errorGet = errors.New("failed to get")
	)

	t.Run("smoke", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(&repository.User{}, nil)
		s.mockRepo.EXPECT().GetSendCoinEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.SendCoinEvent{}, nil)
		s.mockRepo.EXPECT().GetReceiveCoinEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.ReceiveCoinEvent{}, nil)
		s.mockRepo.EXPECT().GetUserItems(gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.Item{}, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		info, err := s.svc.GetUserInfo(ctx, login)

		require.NoError(t, err)
		assert.Equal(t, userInfo, info)
	})

	t.Run("with error getting", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(&repository.User{}, nil)
		s.mockRepo.EXPECT().GetSendCoinEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New(errorGet.Error()))
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Times(1)

		_, err := s.svc.GetUserInfo(ctx, login)

		require.EqualError(t, err, errorGet.Error())
	})
}
