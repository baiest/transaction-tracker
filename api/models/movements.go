package models

import (
	"time"
	"transaction-tracker/internal/movements/domain"
)

type CreateMovementRequest struct {
	InstitutionID string                  `form:"institution_id" `
	Type          domain.MovementType     `form:"type" binding:"required"`
	Amount        float64                 `form:"amount" binding:"required"`
	Date          time.Time               `form:"date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Category      domain.MovementCategory `form:"category" binding:"required"`
	Description   string                  `form:"description"`
	AccountID     string                  `form:"-"`
}

type MovementResponse struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"accountId"`
	InstitutionID  string    `json:"institutionId,omitempty"`
	MessageID      string    `json:"message_id,omitempty"`
	NotificationID string    `json:"notification_id,omitempty"`
	Description    string    `json:"description,omitempty"`
	Amount         float64   `json:"amount"`
	Type           string    `json:"type"`
	Date           time.Time `json:"date"`
	Source         string    `json:"source,omitempty"`
	Category       string    `json:"category,omitempty"`
}

type MovementsListResponse struct {
	TotalPages int64               `json:"total_pages"`
	Page       int64               `json:"page"`
	Movements  []*MovementResponse `json:"movements"`
}

type MovementByYear struct {
	TotalIncome  float64                        `json:"total_income"`
	TotalExpense float64                        `json:"total_expense"`
	Balance      float64                        `json:"balance"`
	Months       []MovementIncomeOutcomeByMonth `json:"months"`
}

type MovementByMonth struct {
	TotalIncome  float64                      `json:"total_income"`
	TotalOutcome float64                      `json:"total_outcome"`
	Year         int                          `json:"year"`
	Balance      float64                      `json:"balance"`
	Days         []MovementIncomeOutcomeByDay `json:"days"`
}

type MovementIncomeOutcome struct {
	Income  float64 `json:"income"`
	Outcome float64 `json:"outcome"`
}

type MovementIncomeOutcomeByMonth struct {
	Month time.Month `json:"month"`
	MovementIncomeOutcome
}
type MovementIncomeOutcomeByDay struct {
	Day int `json:"day"`
	MovementIncomeOutcome
}

func ToDomainMovement(req CreateMovementRequest) *domain.Movement {
	return domain.NewMovement(
		req.AccountID,
		"",
		"",
		"",
		req.Description,
		req.Amount,
		req.Category,
		req.Type,
		req.Date,
		domain.ManualSource,
	)
}

// ToMovementResponse converts a single domain.Movement to an API MovementResponse.
func ToMovementResponse(m *domain.Movement) *MovementResponse {
	return &MovementResponse{
		ID:            m.ID,
		AccountID:     m.AccountID,
		InstitutionID: m.InstitutionID,
		MessageID:     m.MessageID,
		Description:   m.Description,
		Amount:        m.Amount,
		Type:          string(m.Type),
		Date:          m.Date,
		Source:        string(m.Source),
		Category:      string(m.Category),
	}
}

// ToMovementResponses converts a slice of domain.Movement to a slice of API MovementResponse.
func ToMovementResponses(movements []*domain.Movement) []*MovementResponse {
	if movements == nil {
		return []*MovementResponse{}
	}

	response := make([]*MovementResponse, len(movements))
	for i, m := range movements {
		response[i] = ToMovementResponse(m)
	}
	return response
}
