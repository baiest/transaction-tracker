package messageextractor

import (
	"fmt"
	"transaction-tracker/api/services/gmail/models"
)

func NewMovementExtractor(transformerType string, decodedBody string, messageType models.MessageType) (MovementExtractor, error) {
	if messageType == models.Extract {
		return NewDaviviendaExtractor("", messageType), nil
	}

	if decodedBody == "" {
		return nil, fmt.Errorf("missing body")
	}

	return NewDaviviendaExtractor(decodedBody, messageType), nil
}
