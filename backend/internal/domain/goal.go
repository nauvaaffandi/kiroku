package domain

import "time"

type Goal struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TargetAmount float64   `json:"target_amount"`
	Icon         string    `json:"icon,omitempty"`
	DesiredDate  string    `json:"desired_date"` // DD-MM-YYYY
	ColorTheme   string    `json:"color_theme,omitempty"`
	Note         string    `json:"note,omitempty"`
	SavedAlready float64   `json:"saved_already"`
	Currency    string     `json:"currency"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GoalContribution struct {
	ID      string    `json:"id"`
	GoalID  string    `json:"goal_id"`
	Amount  float64   `json:"amount"`
	AddedAt time.Time `json:"added_at"`
	Note    string    `json:"note,omitempty"`
}

type CreateGoalInput struct {
	Name         string  `json:"name"`
	TargetAmount float64 `json:"target_amount"`
	Icon         string  `json:"icon,omitempty"`
	DesiredDate  string  `json:"desired_date,omitempty"` // DD-MM-YYYY from API
	ColorTheme   string  `json:"color_theme,omitempty"`
	Note         string  `json:"note,omitempty"`
	Currency     string  `json:"currency"`
}

type UpdateGoalInput struct {
	Name         *string  `json:"name,omitempty"`
	TargetAmount *float64 `json:"target_amount,omitempty"`
	Icon         *string  `json:"icon,omitempty"`
	DesiredDate  *string  `json:"desired_date,omitempty"` // DD-MM-YYYY from API
	ColorTheme   *string  `json:"color_theme,omitempty"`
	Note         *string  `json:"note,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

type ContributeGoalInput struct {
	Amount float64 `json:"amount"`
	Note   string  `json:"note,omitempty"`
}

type GoalProjection struct {
	CurrentSaved       float64 `json:"current_saved"`
	TargetAmount       float64 `json:"target_amount"`
	Remaining          float64 `json:"remaining"`
	WeeklyRate         float64 `json:"weekly_rate"`
	EstimatedWeeks     int     `json:"estimated_weeks"`
	EstimatedReachDate string  `json:"estimated_reach_date"` // DD-MM-YYYY
	Currency           string  `json:"currency"`
}
