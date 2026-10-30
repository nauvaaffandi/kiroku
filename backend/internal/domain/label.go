package domain

import (
	"context"
	"time"
)

type Label struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon,omitempty"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateLabelInput struct {
	Name        string  `json:"name"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}

type UpdateLabelInput struct {
	Name        *string `json:"name,omitempty"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}

type LabelRepository interface {
	List(ctx context.Context) ([]Label, error)
	Get(ctx context.Context, id string) (*Label, error)
	Create(ctx context.Context, input CreateLabelInput) (*Label, error)
	Update(ctx context.Context, id string, input UpdateLabelInput) (*Label, error)
	Delete(ctx context.Context, id string) error
	Reorder(ctx context.Context, items []ReorderItem) error
}
