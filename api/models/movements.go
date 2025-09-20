package models

import (
	"time"
	"transaction-tracker/internal/movements/domain"
)

type CreateMovementRequest struct {
	InstitutionID string              `form:"institution_id" `
	Type          domain.MovementType `form:"type" binding:"required"`
	Amount        float64             `form:"amount" binding:"required"`
	Date          time.Time           `form:"date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Description   string              `form:"description"`
	AccountID     string              `form:"-"`
}

type MovementResponse struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"accountId"`
	InstitutionID string    `json:"institutionId,omitempty"`
	Description   string    `json:"description,omitempty"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	Date          time.Time `json:"date"`
	Source        string    `json:"source,omitempty"`
	Category      string    `json:"category,omitempty"`
}

type MovementsListResponse struct {
	TotalPages int64               `json:"total_pages"`
	Page       int64               `json:"page"`
	Movements  []*MovementResponse `json:"movements"`
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
