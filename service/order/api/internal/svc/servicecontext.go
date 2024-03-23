package svc

import (
	"github.com/ikun666/k_shop/service/order/api/internal/config"
	"github.com/ikun666/k_shop/service/order/rpc/order"
	"github.com/ikun666/k_shop/service/product/rpc/product"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc    order.Order
	ProductRpc  product.Product
	Redis       *redis.Redis
	KafkaPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		OrderRpc:    order.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		ProductRpc:  product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		Redis:       redis.MustNewRedis(c.RedisConf),
		KafkaPusher: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
