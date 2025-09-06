package transformers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"transaction-tracker/database/mongo/schemas"
)

type Davivienda struct {
	date    time.Time
	hour    time.Time
	value   float64
	details string
	place   string
	text    string
}

var (
	regex = regexp.MustCompile(
		`Fecha:(?P<fecha>.+)\nHora:(?P<hora>.+)\nValor Transacción:\s*(?P<valor>.+)\nClase de Movimiento:\s*(?P<clase>.+),\nLugar de Transacción:(?P<lugar>.+)`)
)

func NewDaviviendaTransformer(text string) MovementTransformer {
	return &Davivienda{
		text: text,
	}
}

func (d *Davivienda) Excecute() (*schemas.Movement, error) {
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

	return schemas.NewMovement("", "", time.Date(
		d.date.Year(), d.date.Month(), d.date.Day(),
		d.hour.Hour(), d.hour.Minute(), d.hour.Second(), 0, d.hour.Location(),
	), d.value, d.value < 0, "others", d.details+" "+d.place), nil
}
