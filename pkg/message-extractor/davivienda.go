package messageextractor

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"transaction-tracker/api/services/gmail/models"
	"transaction-tracker/database/mongo/schemas"
	"transaction-tracker/internal/movements/domain"
	documentextractor "transaction-tracker/pkg/document-extractor"
)

type davivienda struct {
	date        time.Time
	hour        time.Time
	value       float64
	details     string
	place       string
	text        string
	messageType models.MessageType
	extract     *schemas.GmailExtract
}

var (
	daviviendaExtractor = &documentextractor.DaviviendaExtract{
		Password: os.Getenv("EXTRACT_PDF_PASSWORD"),
	}

	regex = regexp.MustCompile(
		`Fecha:(?P<fecha>.+)\nHora:(?P<hora>.+)\nValor Transacción:\s*(?P<valor>.+)\nClase de Movimiento:\s*(?P<clase>.+),\nLugar de Transacción:(?P<lugar>.+)`)
)

func NewDaviviendaExtractor(text string, messageType models.MessageType) MovementExtractor {
	return &davivienda{
		text:        text,
		messageType: messageType,
	}
}

func (d *davivienda) Extract() ([]*domain.Movement, error) {
	switch d.messageType {
	case models.Movement:
		return d.excecuteMovement()
	case models.Extract:
		return d.excecuteExtract()
	}

	return nil, fmt.Errorf("message type not defined")
}

func (d *davivienda) SetExtract(extract *schemas.GmailExtract) {
	d.extract = extract
}

func (d *davivienda) excecuteMovement() ([]*domain.Movement, error) {
	cleanedText := cleanAndNormalizeText(d.text)

	matches := regex.FindStringSubmatch(cleanedText)
	if matches == nil {
		return nil, fmt.Errorf("not found labels: %s", d.text)
	}

	for i, name := range regex.SubexpNames() {
		switch name {
		case "fecha":
			date, err := time.ParseInLocation("2006/01/02", strings.TrimSpace(matches[i]), time.Local)
			if err != nil {
				return nil, fmt.Errorf("error parsing fecha '%s': %v", matches[i], err)
			}

			d.date = date

		case "hora":
			hour, err := time.Parse("15:04:05", strings.TrimSpace(matches[i]))
			if err != nil {
				return nil, fmt.Errorf("error parsing hora '%s': %v", matches[i], err)
			}

			now := time.Now()
			d.hour = time.Date(now.Year(), now.Month(), now.Day(), hour.Hour(), hour.Minute(), hour.Second(), 0, now.Location())

		case "valor":
			clean := strings.ReplaceAll(matches[i], "$", "")
			clean = strings.ReplaceAll(clean, ",", "")
			clean = strings.TrimSpace(clean)

			value, err := strconv.ParseFloat(clean, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing valor '%s': %v", matches[i], err)
			}

			d.value = value

		case "clase":
			d.details = strings.TrimSpace(matches[i])

		case "lugar":
			d.place = strings.TrimSpace(matches[i])
		}
	}

	movementType := domain.Income
	if strings.Contains(d.details, "Descuento") {
		movementType = domain.Expense
	}

	return []*domain.Movement{domain.NewMovement(
		"",
		"",
		"",
		"",
		d.details+" "+d.place,
		d.value,
		domain.Unknown,
		movementType,
		time.Date(
			d.date.Year(), d.date.Month(), d.date.Day(),
			d.hour.Hour(), d.hour.Minute(), d.hour.Second(), 0, d.hour.Location(),
		),
		domain.EmailSource,
	)}, nil
}

func (d *davivienda) excecuteExtract() ([]*domain.Movement, error) {
	if d.extract == nil {
		return nil, fmt.Errorf("missing extract")
	}

	if daviviendaExtractor.Password == "" {
		return nil, fmt.Errorf("missing extract password")
	}

	movementsExtracted := daviviendaExtractor.GetMovements(d.extract.FilePath)
	if len(movementsExtracted) == 0 {
		return nil, fmt.Errorf("missing movements from extract")
	}

	movements := []*domain.Movement{}

	for _, m := range movementsExtracted {
		movementType := domain.Income
		if m.IsNegative {
			movementType = domain.Expense
		}

		movement := domain.NewMovement(
			"",
			"",
			"",
			"",
			m.Detail,
			m.Value,
			domain.Unknown,
			movementType,
			m.Date,
			domain.ExtractSource,
		)

		if movement != nil {
			movements = append(movements, movement)
		}
	}

	return movements, nil
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
