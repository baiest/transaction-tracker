package transformers

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"transaction-tracker/api/services/gmail/models"
	"transaction-tracker/database/mongo/schemas"

	"google.golang.org/api/gmail/v1"
)

type MovementTransformer interface {
	Excecute() ([]*schemas.Movement, error)
	SetExtract(*schemas.GmailExtract)
}

func NewMovementTransformer(transformerType string, msg *gmail.Message, messageType models.MessageType) (MovementTransformer, error) {
	body := ""

	if messageType == models.Extract {
		return NewDaviviendaTransformer(body, messageType), nil
	}

	if len(msg.Payload.Parts) > 0 {
		body = msg.Payload.Parts[0].Body.Data
	} else {
		body = msg.Payload.Body.Data
	}

	if body == "" {
		return nil, fmt.Errorf("missing body")
	}

	decodedBody, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, fmt.Errorf("Error decoding body")
	}

	body = string(decodedBody)

	return NewDaviviendaTransformer(body, messageType), nil
}

func cleanAndNormalizeText(text string) string {
	lines := strings.Split(text, "\n")
	var cleanedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Clase de Movimiento:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				value = regexp.MustCompile(`\s+`).ReplaceAllString(value, " ")

				if !strings.HasSuffix(value, ",") {
					value += ","
				}

				line = "Clase de Movimiento: " + value
			}
		} else if strings.HasPrefix(line, "Valor Transacción:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				line = "Valor Transacción: " + value
			}
		} else if strings.HasPrefix(line, "Lugar de Transacción:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				line = "Lugar de Transacción:" + value
			}
		}

		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
}
