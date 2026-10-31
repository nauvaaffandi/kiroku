package domain

type CalendarDayBackend struct {
	Date             string        `json:"date"`
	Transactions     []CalendarTxn `json:"transactions"`
	Scheduled        []CalendarTxn `json:"scheduled"`
	BalanceDelta     float64       `json:"balance_delta"`
	ProjectedBalance float64       `json:"projected_balance"`
}
