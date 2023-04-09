package registry

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	etcdclient "go.etcd.io/etcd/client/v3"
	"product/internal/conf"
	"time"
)

func InitEtcdService(conf *conf.Bootstrap) error {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints:   []string{conf.Etcd.Url},
		DialTimeout: 5 * time.Second, // 连接超时时间
	})
	if err != nil {
		panic(err)
	}
	conn := etcd.New(client)
	err = conn.Register(context.Background(), &registry.ServiceInstance{
		ID:        "qs-mall-kratos-product",
		Name:      "product",
		Version:   "v1.0.0",
		Endpoints: []string{"grpc://" + conf.Server.Grpc.Addr},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("创建服务成功")
	return err
}
