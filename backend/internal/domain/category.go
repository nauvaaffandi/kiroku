package domain

import "time"

type Category struct {
	ID        string     `json:"id"`
	ParentID  *string    `json:"parent_id"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	Color     string     `json:"color"`
	Kind      *string    `json:"kind"`
	SortOrder int        `json:"sort_order"`
	Depth     int        `json:"depth,omitempty"`
	Path      string     `json:"path,omitempty"`
	Children  []*Category `json:"children,omitempty"`
	Archived  bool       `json:"archived"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateCategoryInput struct {
	Name     string
	ParentID *string
	Icon     string
	Color    string
	Kind     *string
	SortOrder int
}

type UpdateCategoryInput struct {
	Name      *string
	ParentID  *string
	Icon      *string
	Color     *string
	Kind      *string
	Archived  *bool
	SortOrder *int
}

type ReorderItem struct {
	ID        string
	SortOrder int
}
