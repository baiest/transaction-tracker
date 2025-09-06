package models

import "transaction-tracker/database/mongo/schemas"

type MovementsListResponse struct {
	TotalPages int64               `json:"total_pages"`
	Page       int64               `json:"page"`
	Movements  []*schemas.Movement `json:"movements"`
}

type MovementByYear struct {
	TotalIncome  float64                  `json:"total_income"`
	TotalOutcome float64                  `json:"total_outcome"`
	Balance      float64                  `json:"balance"`
	Months       []*MovementIncomeOutcome `json:"months"`
}

type MovementByMonth struct {
	TotalIncome  float64                       `json:"total_income"`
	TotalOutcome float64                       `json:"total_outcome"`
	Year         int                           `json:"year"`
	Balance      float64                       `json:"balance"`
	Days         []*MovementIncomeOutcomeByDay `json:"days"`
}

type MovementIncomeOutcome struct {
	Income  float64 `json:"income"`
	Outcome float64 `json:"outcome"`
}

type MovementIncomeOutcomeByDay struct {
	Day int `json:"day"`
	MovementIncomeOutcome
}
