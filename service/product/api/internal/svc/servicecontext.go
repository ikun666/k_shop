package svc

import (
	"github.com/ikun666/k_shop/service/product/api/internal/config"
	"github.com/ikun666/k_shop/service/product/rpc/product"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	ProductRpc product.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		ProductRpc: product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
