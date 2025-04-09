package etcd

import (
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"my_destributed_project/configs"
	"my_destributed_project/pkg/log"
	"strconv"
	"sync"
	"time"
)

var (
	cli  *clientv3.Client
	once sync.Once
)

// Init 初始化 etcd 客户端
func Init() *clientv3.Client {
	once.Do(func() {
		// 在 Docker 内部，容器可以使用其他容器的名称（即服务名）来进行通信，而无需使用 IP 地址
		etcdPort := strconv.Itoa(configs.AppConfig.Server.EtcdPort)
		endpoints := []string{"http://" + configs.AppConfig.Server.KafkaHost + ":" + etcdPort}

		var err error
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Logger.Error("连接 etcd 失败: ", zap.Error(err))
		}
	})
	return cli
}

// GetClient 获取 etcd 客户端
func GetClient() *clientv3.Client {
	return cli
}
