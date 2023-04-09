package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	v1 "product/api/category/v1"
	"testing"
)

func TestCategoryTreeListService(t *testing.T) {
	// 连接 gRPC 服务
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := v1.NewCategoryClient(conn)

	// 调用远程方法
	_, err = client.CategoryTreeList(context.Background(), &v1.CategoryListRequest{})
	if err != nil {
		log.Fatalf("failed to call remote method: %v", err)
	}

	// 处理响应结果
	fmt.Println("test")
}
