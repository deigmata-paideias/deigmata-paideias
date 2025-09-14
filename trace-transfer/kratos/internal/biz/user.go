package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type IUserBiz interface {
	GetUser(ctx context.Context, id string) (string, error)
}

type IUserRepo interface {
	GetUser(ctx context.Context, id string) (string, error)
}

type UserBiz struct {
	userRepo IUserRepo
	log      *log.Helper
}

func NewUserBiz(userRepo IUserRepo, logger log.Logger) IUserBiz {

	return &UserBiz{
		userRepo: userRepo,
		log:      log.NewHelper(logger),
	}
}

func (b *UserBiz) GetUser(ctx context.Context, id string) (string, error) {

	user, err := b.userRepo.GetUser(ctx, id)
	if err != nil {
		b.log.WithContext(ctx).Error("Failed to get user", "id", id, "error", err)
		return "", err
	}

	b.log.WithContext(ctx).Info("Successfully retrieved user", "id", id, "user", user)

	return user, nil
}
