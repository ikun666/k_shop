package svc

import (
	"github.com/ikun666/k_shop/service/order/mq/internal/config"
	"github.com/ikun666/k_shop/service/order/rpc/order"
	"github.com/ikun666/k_shop/service/product/rpc/product"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	ProductRPC product.Product
	OrderRPC   order.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		ProductRPC: product.NewProduct(zrpc.MustNewClient(c.ProductRPC)),
		OrderRPC:   order.NewOrder(zrpc.MustNewClient(c.OrderRPC)),
	}

}
