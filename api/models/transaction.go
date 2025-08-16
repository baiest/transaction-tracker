package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TransactionRequest struct {
	Date  time.Time `json:"date"`
	Type  string    `json:"type"`
	Value uint64    `json:"value"`
}

type DaviviendaEmailMovement struct {
	Date  time.Time
	Hour  time.Time
	Value float64
	Type  string
	Place string
}

var (
	regex = regexp.MustCompile(
		`Fecha:(?P<fecha>.+)\nHora:(?P<hora>.+)\nValor Transacción:\s*(?P<valor>.+)\nClase de Movimiento:\s*(?P<clase>.+),\nLugar de Transacción:(?P<lugar>.+)`)
)

func (dm *DaviviendaEmailMovement) ToTransactionRequest() *TransactionRequest {
	tr := &TransactionRequest{
		Type:  "others",
		Value: uint64(dm.Value),
	}

	tr.Date = time.Date(
		dm.Date.Year(), dm.Date.Month(), dm.Date.Day(),
		dm.Hour.Hour(), dm.Hour.Minute(), dm.Hour.Second(), 0, dm.Date.Location(),
	)

	return tr
}

func NewDaviviendaMovementFromText(text string) (*DaviviendaEmailMovement, error) {
	cleanedText := cleanAndNormalizeText(text)

	matches := regex.FindStringSubmatch(cleanedText)
	if matches == nil {
		return nil, fmt.Errorf("not found labels: %s", text)
	}

	result := &DaviviendaEmailMovement{}
	for i, name := range regex.SubexpNames() {
		switch name {
		case "fecha":
			date, err := time.ParseInLocation("2006/01/02", strings.TrimSpace(matches[i]), time.Local)
			if err != nil {
				return nil, fmt.Errorf("error parsing fecha '%s': %v", matches[i], err)
			}

			result.Date = date

		case "hora":
			hour, err := time.Parse("15:04:05", strings.TrimSpace(matches[i]))
			if err != nil {
				return nil, fmt.Errorf("error parsing hora '%s': %v", matches[i], err)
			}

			now := time.Now()
			result.Hour = time.Date(now.Year(), now.Month(), now.Day(), hour.Hour(), hour.Minute(), hour.Second(), 0, now.Location())

		case "valor":
			clean := strings.ReplaceAll(matches[i], "$", "")
			clean = strings.ReplaceAll(clean, ",", "")
			clean = strings.TrimSpace(clean)

			value, err := strconv.ParseFloat(clean, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing valor '%s': %v", matches[i], err)
			}

			result.Value = value

		case "clase":
			result.Type = strings.TrimSpace(matches[i])

		case "lugar":
			result.Place = strings.TrimSpace(matches[i])
		}
	}

	return result, nil
}

// cleanAndNormalizeText cleans and normalizes the input text to match the expected format
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

func (tr *TransactionRequest) MarshalJSON() ([]byte, error) {
	type Alias TransactionRequest

	aux := &struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(tr),
		Date:  tr.Date.UTC().Format("2006-01-02T15:04:05Z"),
	}

	return json.Marshal(aux)
}
