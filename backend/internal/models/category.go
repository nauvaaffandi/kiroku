package models

type Category struct {
	ID        string     `json:"id"`
	ParentID  *string    `json:"parent_id,omitempty"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon,omitempty"`
	Color     string     `json:"color,omitempty"`
	Kind      *string    `json:"kind,omitempty"`
	SortOrder int        `json:"sort_order"`
	Depth     int        `json:"depth"`
	Path      string     `json:"path,omitempty"`
	Children  []Category `json:"children,omitempty"`
}

type CreateCategoryReq struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
	Icon     string  `json:"icon,omitempty"`
	Color    string  `json:"color,omitempty"`
	Kind     *string `json:"kind,omitempty"`
}

type UpdateCategoryReq struct {
	Name     *string `json:"name,omitempty"`
	ParentID *string `json:"parent_id,omitempty"`
	Icon     *string `json:"icon,omitempty"`
	Color    *string `json:"color,omitempty"`
	Kind     *string `json:"kind,omitempty"`
	Archived *bool   `json:"archived,omitempty"`
}

type ReorderReq struct {
	Items []struct {
		ID        string `json:"id"`
		SortOrder int    `json:"sort_order"`
	} `json:"items"`
}
