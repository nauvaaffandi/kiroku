package domain

import "context"

type CategoryRepository interface {
	List(ctx context.Context, includeArchived bool) ([]Category, error)
	Tree(ctx context.Context) ([]Category, error)
	Get(ctx context.Context, id string) (*Category, error)
	Create(ctx context.Context, input CreateCategoryInput) (*Category, error)
	Update(ctx context.Context, id string, input UpdateCategoryInput) (*Category, error)
	Delete(ctx context.Context, id string) error
	Reorder(ctx context.Context, items []ReorderItem) error
	CountChildren(ctx context.Context, id string) (int, error)
	NextSortOrder(ctx context.Context, parentID *string) (int, error)
}
