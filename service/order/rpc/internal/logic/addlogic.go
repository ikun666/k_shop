package logic

import (
	"context"

	"github.com/ikun666/k_shop/service/order/rpc/internal/svc"
	"github.com/ikun666/k_shop/service/order/rpc/types/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddLogic) Add(in *order.CreateRequest) (*order.CreateResponse, error) {
	// todo: add your logic here and delete this line

	return &order.CreateResponse{}, nil
}
