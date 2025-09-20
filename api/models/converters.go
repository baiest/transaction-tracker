package models

import "transaction-tracker/internal/movements/domain"

func ToDomainMovement(req CreateMovementRequest) *domain.Movement {
	return &domain.Movement{
		AccountID:     req.AccountID,
		InstitutionID: req.InstitutionID,
		Type:          req.Type,
		Amount:        req.Amount,
		Date:          req.Date,
		Description:   req.Description,
	}
}

// ToMovementResponse converts a single domain.Movement to an API MovementResponse.
func ToMovementResponse(m *domain.Movement) *MovementResponse {
	return &MovementResponse{
		ID:            m.ID,
		AccountID:     m.AccountID,
		InstitutionID: m.InstitutionID,
		Description:   m.Description,
		Amount:        m.Amount,
		Type:          string(m.Type),
		Date:          m.Date,
		Source:        string(m.Source),
		Category:      m.Category,
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
