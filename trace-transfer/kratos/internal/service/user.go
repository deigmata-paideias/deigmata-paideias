package service

import (
	"context"

	v1 "kratos/api/user/v1"
	"kratos/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

type UserService struct {
	v1.UnimplementedUserServer
	log *log.Helper

	biz biz.IUserBiz
}

func NewUserService(biz biz.IUserBiz, logger log.Logger) *UserService {

	return &UserService{
		biz: biz,
		log: log.NewHelper(logger),
	}
}

// GetUser implements v1.UserServer
func (s *UserService) GetUser(ctx context.Context, req *v1.UserRequest) (*v1.UserResponse, error) {

	// 打印所有 header 便于链路排查
	tr, ok := transport.FromServerContext(ctx)
	if ok {
		headers := tr.RequestHeader()
		for _, k := range headers.Keys() {
			s.log.WithContext(ctx).Infof("[kratos] header: %s = %s", k, headers.Get(k))
		}
	}

	user, err := s.biz.GetUser(ctx, req.GetId())
	if err != nil {
		s.log.WithContext(ctx).Infof("GetUser err: %v", err)
		return nil, err
	}

	return &v1.UserResponse{
		Name: user,
	}, nil
}
