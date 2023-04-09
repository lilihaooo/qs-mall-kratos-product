package service

import (
	"context"
	v1 "product/api/category/v1"
	"product/internal/biz"
)

// GreeterService is a greeter service.
type CategoryService struct {
	v1.UnimplementedCategoryServer
	uc *biz.CategoryUsecase
}

// NewGreeterService new a greeter service.
func NewCategoryService(uc *biz.CategoryUsecase) *CategoryService {
	return &CategoryService{uc: uc}
}

func (s *CategoryService) CategoryCreate(ctx context.Context, in *v1.CategoryCreateRequest) (*v1.Reply, error) {
	err := s.uc.CreateCategory(ctx, &biz.Category{Name: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.Reply{Code: 20001, Message: "添加成功"}, nil
}

// CategoryTreeList RPC方法
func (s *CategoryService) CategoryTreeList(ctx context.Context, in *v1.CategoryListRequest) (*v1.CategoryListReply, error) {
	res, err := s.uc.CategoryTreeList(ctx)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
