package data

import (
	"context"

	"kratos-example/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type exampleRepo struct {
	data *Data
	log  *log.Helper
}

// NewExampleRepo .
func NewExampleRepo(data *Data, logger log.Logger) biz.ExampleRepo {
	return &exampleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *exampleRepo) Save(ctx context.Context, g *biz.Example) (*biz.Example, error) {
	return g, nil
}

func (r *exampleRepo) Update(ctx context.Context, g *biz.Example) (*biz.Example, error) {
	return g, nil
}

func (r *exampleRepo) FindByID(context.Context, int64) (*biz.Example, error) {
	return nil, nil
}

func (r *exampleRepo) ListByHello(context.Context, string) ([]*biz.Example, error) {
	return nil, nil
}

func (r *exampleRepo) ListAll(context.Context) ([]*biz.Example, error) {
	return nil, nil
}
