package transport

import (
	"context"
	"fmt"

	clrServer "collector-go/internal/collector/common/server"
	"collector-go/internal/collector/common/types"
)

type Config struct {
	clrServer.Server
}

type Runner struct {
	Config
}

func New(srv *Config) *Runner {
	return &Runner{
		Config: *srv,
	}
}

func (r *Runner) Start(ctx context.Context) error {
	fmt.Println("transport start...")

	select {
	case <-ctx.Done():
		return nil
	}
}

func (r *Runner) Info() types.Info {
	return types.Info{
		Name: "transport",
	}
}

func (r *Runner) Close() error {
	fmt.Println("transport close...")
	return nil
}
