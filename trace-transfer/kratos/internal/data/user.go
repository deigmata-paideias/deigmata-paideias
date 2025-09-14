package data

import (
	"context"
	"kratos/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.IUserRepo {

	return &UserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *UserRepo) GetUser(ctx context.Context, id string) (string, error) {

	data := r.data.db.userDB
	user, exists := data[id]

	if !exists {

		r.log.WithContext(ctx).Error("User not found. ", "id", id)
		return "", nil
	}
	r.log.WithContext(ctx).Info("User retrieved successfully. ", "id", id)

	return user, nil
}
