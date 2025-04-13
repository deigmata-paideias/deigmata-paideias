package service

import (
	"context"

	v1 "kratos-example/api/example/v1"
	"kratos-example/internal/biz"
)

// ExampleService is a example service.
type ExampleService struct {
	v1.UnimplementedExampleServer

	uc *biz.ExampleUsecase
}

// NewExampleService new a greeter service.
func NewExampleService(uc *biz.ExampleUsecase) *ExampleService {
	return &ExampleService{uc: uc}
}

// SayHello implements example.ExampleServer.
func (s *ExampleService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateExample(ctx, &biz.Example{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
