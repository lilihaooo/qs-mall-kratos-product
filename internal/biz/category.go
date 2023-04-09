package biz

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"

	v1 "product/api/category/v1"

	"github.com/go-kratos/kratos/v2/log"
)

var (
// ErrUserNotFound is user not found.
// ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Category is a Category model.
type Category struct {
	gorm.Model
	Name         string
	ParentCid    int64
	CatLevel     int32
	ShowStatus   int32
	Sort         int32
	Icon         string
	ProductUnit  string
	ProductCount int64
}

// TableName 指定表名
func (Category) TableName() string {
	return "pms_category"
}

// 自己嵌套自己, 同事返回数组就可以实现无限分类

type CategoryTree struct {
	ID           int64
	Name         string
	ParentCid    int64
	CatLevel     int32
	ShowStatus   int32
	Sort         int32
	Icon         string
	ProductUnit  string
	ProductCount int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Children     []CategoryTree
}

// CategoryRepo is a Greater repo.
type CategoryRepo interface {
	List(context.Context) ([]CategoryTree, error)
	Create(context.Context, *Category) error
	Delete(context.Context, int64) error
	Update(context.Context, int64, *Category) error
}

// CategoryUsecase is a Category usecase.
type CategoryUsecase struct {
	repo CategoryRepo
	log  *log.Helper
}

// NewCategoryUsecase new a Category usecase.
func NewCategoryUsecase(repo CategoryRepo, logger log.Logger) *CategoryUsecase {
	return &CategoryUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateCategory creates a Category, and returns the new Category.
func (uc *CategoryUsecase) CreateCategory(ctx context.Context, g *Category) error {
	uc.log.WithContext(ctx).Infof("CreateCategory: %v", g.Name)
	return uc.repo.Create(ctx, g)
}

// CategoryTreeList 业务逻辑
func (uc *CategoryUsecase) CategoryTreeList(ctx context.Context) (v1.CategoryListReply, error) {
	CategoryTreeRes, err := uc.repo.List(ctx)
	if err != nil {
		fmt.Println("数据库出错")
		fmt.Println(err.Error())
	}
	// 使用递归将查询结果类型:[]CategoryTree, 转为[]*v1.CategoryTree
	pbCategoryTree := CategoryTreeToV1CategoryTree(CategoryTreeRes)

	// 找儿子
	pbRootCategories := make([]*v1.CategoryTree, 0, len(pbCategoryTree))
	for _, c := range pbCategoryTree {
		if c.ParentCid == 0 {
			pbRootCategories = append(pbRootCategories, c)
		} else {
			for _, parent := range pbCategoryTree {
				if parent.Id == c.ParentCid {
					parent.Children = append(parent.Children, c)
					break
				}
			}
		}
	}
	// 将转换后的消息结构体返回
	return v1.CategoryListReply{Code: 20001, Message: "获取分类树成功", Data: pbRootCategories}, nil
}

func CategoryTreeToV1CategoryTree(CategoryTree []CategoryTree) []*v1.CategoryTree {
	v1CategoryTree := make([]*v1.CategoryTree, 0, len(CategoryTree))
	for _, c := range CategoryTree {
		v1CategoryTree = append(v1CategoryTree, &v1.CategoryTree{
			Id:        c.ID,
			Name:      c.Name,
			ParentCid: c.ParentCid,
			CatLevel:  c.CatLevel,
			Children:  CategoryTreeToV1CategoryTree(c.Children),
		})
	}
	return v1CategoryTree
}
