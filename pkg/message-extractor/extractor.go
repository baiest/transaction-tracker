package messageextractor

import (
	"fmt"
)

func NewMovementExtractor(transformerType string, decodedBody string, messageType MessageType) (MovementExtractor, error) {
	if messageType == Movement && decodedBody == "" {
		return nil, fmt.Errorf("missing body")
	}

	return NewDaviviendaExtractor(decodedBody, messageType), nil
}
