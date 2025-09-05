package job

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
	fmt.Println("job start...")

	select {
	case <-ctx.Done():
		return nil
	}
}

func (r *Runner) Info() types.Info {
	return types.Info{
		Name: "job",
	}
}

func (r *Runner) Close() error {
	fmt.Println("job close...")
	return nil
}
