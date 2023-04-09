package data

import (
	"context"
	"gorm.io/gorm"
	"product/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	CAT_LEVEL_ONE   = 1
	CAT_LEVEL_TWO   = 2
	CAT_LEVEL_THREE = 3
)

type CategoryRepo struct {
	data *Data
	log  *log.Helper
}

func NewCategoryRepo(data *Data, logger log.Logger) biz.CategoryRepo {
	return &CategoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *CategoryRepo) List(ctx context.Context) ([]biz.CategoryTree, error) {
	var categories []*biz.Category
	if err := r.data.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return GetCategoryTree(categories, 0)
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *CategoryRepo) Update(ctx context.Context, id int64, c *biz.Category) error {
	return nil
}

func (r *CategoryRepo) Create(ctx context.Context, c *biz.Category) error {
	return nil
}

func GetCategoryTree(categories []*biz.Category, parentCID int64) ([]biz.CategoryTree, error) {
	var categoryTrees []biz.CategoryTree

	for _, category := range categories {
		if category.ParentCid == parentCID {
			categoryTree := biz.CategoryTree{
				ID:        category.ID,
				Name:      category.Name,
				ParentCid: category.ParentCid,
			}
			subcategories, err := GetCategoryTree(categories, category.ID)
			if err != nil {
				return nil, err
			}
			categoryTree.Children = subcategories
			categoryTrees = append(categoryTrees, categoryTree)
		}
	}
	return categoryTrees, nil

}

func GetCategoryTree0(db *gorm.DB, parentCID int64) ([]biz.CategoryTree, error) {
	var categories []*biz.Category
	var categoryTrees []biz.CategoryTree

	if err := db.Where("parent_cid = ?", parentCID).Find(&categories).Error; err != nil {
		return nil, err
	}

	for _, category := range categories {
		categoryTree := biz.CategoryTree{
			ID:        category.ID,
			Name:      category.Name,
			ParentCid: category.ParentCid,
		}

		subcategories, err := GetCategoryTree0(db, category.ID)
		if err != nil {
			return nil, err
		}
		categoryTree.Children = subcategories
		categoryTrees = append(categoryTrees, categoryTree)
	}

	return categoryTrees, nil
}
