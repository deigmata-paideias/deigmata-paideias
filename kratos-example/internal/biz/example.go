package biz

import (
	"context"

	v1 "kratos-example/api/example/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Example is a Greeter model.
type Example struct {
	Hello string
}

// ExampleRepo is a Greater repo.
type ExampleRepo interface {
	Save(context.Context, *Example) (*Example, error)
	Update(context.Context, *Example) (*Example, error)
	FindByID(context.Context, int64) (*Example, error)
	ListByHello(context.Context, string) ([]*Example, error)
	ListAll(context.Context) ([]*Example, error)
}

// ExampleUsecase is a Example usecase.
type ExampleUsecase struct {
	repo ExampleRepo
	log  *log.Helper
}

// NewExampleUsecase new a Greeter usecase.
func NewExampleUsecase(repo ExampleRepo, logger log.Logger) *ExampleUsecase {
	return &ExampleUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateExample creates a example, and returns the new Example.
func (uc *ExampleUsecase) CreateExample(ctx context.Context, g *Example) (*Example, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
