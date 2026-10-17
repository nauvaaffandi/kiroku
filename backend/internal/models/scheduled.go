package models

import (
	"encoding/json"
	"time"
)

type Scheduled struct {
	ID                      string          `json:"id"`
	RRule                   string          `json:"rrule"`
	MaxOccurrences          *int            `json:"max_occurrences,omitempty"`
	TransactionTemplateJSON json.RawMessage `json:"transaction_template_json"`
	Enabled                 bool            `json:"enabled"`
	NextRun                 *time.Time      `json:"next_run,omitempty"`
	OccurrencesDone         int             `json:"occurrences_done"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

type CreateScheduledReq struct {
	RRule                   string          `json:"rrule"`
	MaxOccurrences          *int            `json:"max_occurrences,omitempty"`
	TransactionTemplateJSON json.RawMessage `json:"transaction_template_json"`
	Enabled                 *bool           `json:"enabled,omitempty"`
}

type UpdateScheduledReq struct {
	RRule                   *string         `json:"rrule,omitempty"`
	MaxOccurrences          *int            `json:"max_occurrences,omitempty"`
	TransactionTemplateJSON json.RawMessage `json:"transaction_template_json,omitempty"`
	Enabled                 *bool           `json:"enabled,omitempty"`
}

type PresetReq struct {
	Preset   string                 `json:"preset"`
	Params   map[string]interface{} `json:"params"`
	Template map[string]interface{} `json:"template"`
}
