package categories

import (
	"database/sql"

	"kiroku/backend/internal/models"
)

func (h *Handler) computeDepthAndPath(c *models.Category, catMap map[string]*models.Category) {
	if c.ParentID == nil {
		c.Depth = 0
		c.Path = c.Name
		return
	}
	parent, ok := catMap[*c.ParentID]
	if !ok {
		c.Depth = 0
		c.Path = c.Name
		return
	}
	if parent.Depth == 0 && parent.Path == "" {
		h.computeDepthAndPath(parent, catMap)
	}
	c.Depth = parent.Depth + 1
	c.Path = parent.Path + " > " + c.Name
}

func (h *Handler) isDescendant(ancestor, potentialDescendant string) bool {
	current := potentialDescendant
	for {
		var parentID sql.NullString
		err := h.db.QueryRow("SELECT parent_id FROM categories WHERE id = ?", current).Scan(&parentID)
		if err != nil || !parentID.Valid {
			return false
		}
		if parentID.String == ancestor {
			return true
		}
		current = parentID.String
	}
}
