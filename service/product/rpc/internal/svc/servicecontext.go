package svc

import (
	"github.com/ikun666/k_shop/service/product/model"
	"github.com/ikun666/k_shop/service/product/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	ProductModel model.ProductModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:       c,
		ProductModel: model.NewProductModel(conn, c.CacheRedis),
	}
}
