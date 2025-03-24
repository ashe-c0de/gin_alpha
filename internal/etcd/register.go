package etcd

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"my_destributed_project/pkg/log"
)

// RegisterService 注册服务
func RegisterService(serviceKey, serviceValue string, ttl int64) {
	cli := GetClient()
	if cli == nil {
		log.Logger.Error("Etcd 客户端未初始化")
		return
	}

	// 创建租约
	resp, err := cli.Grant(context.TODO(), ttl)
	if err != nil {
		log.Logger.Error("创建租约失败: ", zap.Error(err))
		return
	}

	// 注册服务
	_, err = cli.Put(context.TODO(), serviceKey, serviceValue, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Logger.Error("注册服务失败: ", zap.Error(err))
		return
	}

	// 保持租约
	keepAliveCh, err := cli.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		log.Logger.Error("保持租约失败: ", zap.Error(err))
		return
	}
	log.Logger.Info("4")
	// 监听 keepAlive 响应（防止 `nil` panic）
	go func() {
		for ka := range keepAliveCh {
			if ka == nil {
				log.Logger.Warn("KeepAlive 信道关闭")
				return
			}
		}
	}()

	log.Logger.Info("已注册到 etcd", zap.String("serviceKey", serviceKey))
}
