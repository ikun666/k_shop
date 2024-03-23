package mqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/ikun666/k_shop/service/order/mq/internal/svc"
	"github.com/ikun666/k_shop/service/order/rpc/order"
	"github.com/ikun666/k_shop/service/product/rpc/product"
)

type PaymentSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}
type KafkaData struct {
	Uid    int64 `json:"uid"`
	Pid    int64 `json:"pid"`
	Amount int64 `json:"amount"`
	Status int64 `json:"status"`
}

func NewPaymentSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentSuccess {
	return &PaymentSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 消费数据 实现kafaka接口
func (l *PaymentSuccess) Consume(key, val string) error {
	fmt.Printf("PaymentSuccess key :%s , val :%s\n", key, val)
	msg := &KafkaData{}
	err := json.Unmarshal([]byte(val), msg)
	if err != nil {
		return err
	}
	// 获取 OrderRpc BuildTarget
	orderRpcBusiServer, err := l.svcCtx.Config.OrderRPC.BuildTarget()
	if err != nil {
		return nil
	}

	// 获取 ProductRpc BuildTarget
	productRpcBusiServer, err := l.svcCtx.Config.ProductRPC.BuildTarget()
	if err != nil {
		return nil
	}

	// dtm 服务的 etcd 注册地址
	var dtmServer = "etcd://192.168.44.132:20000/dtmservice"
	// 创建一个gid
	gid := dtmgrpc.MustGenGid(dtmServer)
	// 创建一个saga协议的事务
	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(orderRpcBusiServer+"/order.Order/Create", orderRpcBusiServer+"/order.Order/CreateRevert", &order.CreateRequest{
			Uid:    msg.Uid,
			Pid:    msg.Pid,
			Amount: msg.Amount,
			Status: 0,
		}).
		Add(productRpcBusiServer+"/product.Product/DecrStock", productRpcBusiServer+"/product.Product/DecrStockRevert", &product.DecrStockRequest{
			Id:  msg.Pid,
			Num: 1,
		})

	// 事务提交
	err = saga.Submit()
	if err != nil {
		return nil
	}

	return nil
}
