package transformers

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"transaction-tracker/api/services/gmail/models"
	"transaction-tracker/database/mongo/schemas"
	documentextractor "transaction-tracker/document-extractor"
)

type Davivienda struct {
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

func NewDaviviendaTransformer(text string, messageType models.MessageType) MovementTransformer {
	return &Davivienda{
		text:        text,
		messageType: messageType,
	}
}

func (d *Davivienda) Excecute() ([]*schemas.Movement, error) {
	switch d.messageType {
	case models.Movement:
		return d.excecuteMovement()
	case models.Extract:
		return d.excecuteExtract()
	}

	return nil, fmt.Errorf("message type not defined")
}

func (d *Davivienda) SetExtract(extract *schemas.GmailExtract) {
	d.extract = extract
}

func (d *Davivienda) excecuteMovement() ([]*schemas.Movement, error) {
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

	return []*schemas.Movement{schemas.NewMovement("", "", time.Date(
		d.date.Year(), d.date.Month(), d.date.Day(),
		d.hour.Hour(), d.hour.Minute(), d.hour.Second(), 0, d.hour.Location(),
	), d.value, strings.Contains(d.details, "Descuento"), "others", d.details+" "+d.place)}, nil
}

func (d *Davivienda) excecuteExtract() ([]*schemas.Movement, error) {
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

	movements := []*schemas.Movement{}

	for _, m := range movementsExtracted {
		movement := schemas.NewMovement(
			d.extract.Email,
			d.extract.ID,
			m.Date,
			m.Value,
			m.IsNegative,
			m.Topic,
			m.Detail,
		)

		if movement != nil {
			movements = append(movements, movement)
		}
	}

	return movements, nil
}
