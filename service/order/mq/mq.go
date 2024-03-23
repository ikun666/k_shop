package main

import (
	"context"
	"flag"

	"github.com/ikun666/k_shop/service/order/mq/internal/config"
	"github.com/ikun666/k_shop/service/order/mq/internal/svc"
	"github.com/ikun666/k_shop/service/order/mq/mqs"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/mq.yaml", "the etc file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// server := rest.MustNewServer(c.RestConf)
	// defer server.Stop()
	logx.SetLevel(logx.ErrorLevel)
	logx.DisableStat()
	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range mqs.Consumers(c, ctx, svcCtx) {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}
