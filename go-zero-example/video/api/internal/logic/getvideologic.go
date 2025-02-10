package logic

import (
	"context"
	"fmt"
	"go-zero-example/user/rpc/types/user"

	"go-zero-example/video/api/internal/svc"
	"go-zero-example/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoLogic {
	return &GetVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVideoLogic) GetVideo(req *types.VideoReq) (resp *types.VideoRes, err error) {

	userl, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.IdRequest{
		Id: req.Id,
	})

	fmt.Printf("req.Id: %v, userl.Id: %v\n", req.Id, userl.Id)

	if err != nil {
		return nil, err
	}

	return &types.VideoRes{
		Id:   userl.Id,
		Name: userl.Name,
	}, nil
}
