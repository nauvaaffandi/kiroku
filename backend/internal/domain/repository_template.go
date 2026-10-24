package domain

import "context"

type TemplateRepository interface {
	List(ctx context.Context) ([]Template, error)
	Get(ctx context.Context, id string) (*Template, error)
	Create(ctx context.Context, input CreateTemplateInput) (*Template, error)
	Update(ctx context.Context, id string, input UpdateTemplateInput) (*Template, error)
	Delete(ctx context.Context, id string) error
}
