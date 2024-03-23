package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ikun666/k_shop/service/order/api/internal/svc"
	"github.com/ikun666/k_shop/service/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}
type KafkaData struct {
	Uid    int64 `json:"uid"`
	Pid    int64 `json:"pid"`
	Amount int64 `json:"amount"`
	Status int64 `json:"status"`
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 提前在redis设置库存 hset product:pid total 1000   hset product_pid seckill 0
const (
	luaCheckAndUpdateScript = `
  local counts = redis.call("HMGET", KEYS[1], "total", "seckill")
  local total = tonumber(counts[1])
  local seckill = tonumber(counts[2])
  if seckill + 1 <= total then
	redis.call("HINCRBY", KEYS[1], "seckill", 1)
	return 1
  end
  return 0
  `
)

func (l *CreateLogic) Create(req *types.CreateRequest) (resp *types.CreateResponse, err error) {
	//etcd 分布式锁
	// locker := etcdlock.NewLocker(etcdlock.NewSession())
	// locker.Lock()
	// defer locker.Unlock()
	//redis更新库存

	// value, err := l.svcCtx.Redis.Get(fmt.Sprintf("product_%d", req.Pid))
	// if err != nil {
	// 	return nil, err
	// }
	// if len(value) == 0 {
	// 	res, err := l.svcCtx.ProductRpc.Detail(l.ctx, &product.DetailRequest{
	// 		Id: req.Pid,
	// 	})
	// 	if err != nil {
	// 		if err == model.ErrNotFound {
	// 			return nil, status.Error(100, "产品不存在")
	// 		}
	// 		return nil, status.Error(500, err.Error())
	// 	}
	// 	value = strconv.Itoa(int(res.Stock))

	// 	err = l.svcCtx.Redis.Set(fmt.Sprintf("product_%d", req.Pid), value)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// value_int, err := strconv.Atoi(value)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("stock:", value_int)
	// if value_int > 0 {
	// 	value = strconv.Itoa(value_int - 1)
	// 	// fmt.Println("参数", value_int, in.Count, value)
	// 	err = l.svcCtx.Redis.Set(fmt.Sprintf("product_%d", req.Pid), value)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	//发送到消息队列 数据库更新库存
	// 	msg, err := json.Marshal(KafkaData{
	// 		Uid:    req.Uid,
	// 		Pid:    req.Pid,
	// 		Amount: req.Amount,
	// 		Status: req.Status,
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if err := l.svcCtx.KafkaPusher.Push(string(msg)); err != nil {
	// 		return nil, err
	// 	}
	// 	return &types.CreateResponse{}, nil
	// }
	// return &types.CreateResponse{}, status.Error(500, "库存不足")
	//redis 分布式锁 faster
	value, err := l.svcCtx.Redis.EvalCtx(l.ctx, luaCheckAndUpdateScript, []string{fmt.Sprintf("product:%d", req.Pid)})
	if err != nil {
		return nil, err
	}
	if value.(int64) == 0 {
		return &types.CreateResponse{}, status.Error(500, "库存不足")
	}
	fmt.Println(value)
	//发送到消息队列 数据库更新库存
	msg, err := json.Marshal(KafkaData{
		Uid:    req.Uid,
		Pid:    req.Pid,
		Amount: req.Amount,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.KafkaPusher.Push(string(msg)); err != nil {
		return nil, err
	}
	return &types.CreateResponse{}, nil
}
