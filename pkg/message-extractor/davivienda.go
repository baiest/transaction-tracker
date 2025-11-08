package messageextractor

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	extractDomain "transaction-tracker/internal/extracts/domain"
	movementDomain "transaction-tracker/internal/movements/domain"
	documentextractor "transaction-tracker/pkg/document-extractor"
	"unicode/utf8"
)

const institutionID = "davivienda"

var (
	extractTextFromPDF = documentextractor.ExtractTextFromPDF
	movementRegex      = regexp.MustCompile(`(?m)^(\d{2}\s+\d{2})\s+\$\s*([\d,]+\.\d{2})([+-])\s+(\d{4})\s+(.+)$`)
	yearRegex          = regexp.MustCompile(`INFORME DEL MES:.*?/(\d{4})`)
	regex              = regexp.MustCompile(
		`Fecha\s*:\s*(?P<fecha>[0-9]{4}[/-][0-9]{2}[/-][0-9]{2})\s*Hora\s*:\s*(?P<hora>[0-9:]+)\s*Valor\s+Transacci(?:ón|on)\s*:\s*\$?\s*(?P<valor>[0-9.,]+)\s*Clase\s+de\s+Movimiento\s*:\s*(?P<clase>[^,.;\n]+)[,.;\s]*Lugar\s+de\s+Transacci(?:ón|on)\s*:\s*(?P<lugar>.+)`)
)

type davivienda struct {
	date        time.Time
	hour        time.Time
	value       float64
	details     string
	place       string
	text        string
	messageType MessageType
	extract     *extractDomain.Extract
	password    string
}

func NewDaviviendaExtractor(text string, messageType MessageType) MovementExtractor {
	return &davivienda{
		text:        text,
		messageType: messageType,
		password:    os.Getenv("EXTRACT_PDF_PASSWORD"),
	}
}

func (d *davivienda) Extract() ([]*movementDomain.Movement, error) {
	switch d.messageType {
	case Movement:
		return d.excecuteMovement()
	case Extract:
		return d.excecuteExtract()
	case Unknown:
		return []*movementDomain.Movement{}, nil
	}

	return nil, fmt.Errorf("message type not defined: %s", d.messageType)
}

func (d *davivienda) SetExtract(extract *extractDomain.Extract) {
	d.extract = extract
}

func (d *davivienda) excecuteMovement() ([]*movementDomain.Movement, error) {
	cleanedText := cleanAndNormalizeText(d.text)

	matches := regex.FindStringSubmatch(cleanedText)
	if matches == nil {
		return nil, fmt.Errorf("not found labels: %s", cleanedText)
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

	movementType := movementDomain.Income
	if strings.Contains(d.details, "Descuento") {
		movementType = movementDomain.Expense
	}

	return []*movementDomain.Movement{movementDomain.NewMovement(
		"",
		institutionID,
		"",
		"",
		d.details+" "+d.place,
		d.value,
		movementDomain.Unknown,
		movementType,
		time.Date(
			d.date.Year(), d.date.Month(), d.date.Day(),
			d.hour.Hour(), d.hour.Minute(), d.hour.Second(), 0, d.hour.Location(),
		),
		movementDomain.EmailSource,
	)}, nil
}

func (d *davivienda) excecuteExtract() ([]*movementDomain.Movement, error) {
	if d.extract == nil {
		return nil, fmt.Errorf("missing extract")
	}

	if d.password == "" {
		return nil, fmt.Errorf("missing extract password")
	}

	movements, err := d.GetMovements()
	if err != nil {
		return nil, fmt.Errorf("error extracting movements: %v", err)
	}

	return movements, nil
}

func cleanAndNormalizeText(text string) string {
	// Normalizar retornos de carro y asegurar UTF-8 válido
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = ToValidUTF8(text)

	// Corregir mojibake/encodings comunes que aparecen en PDFs
	text = strings.ReplaceAll(text, "TransacciÃ³n", "Transacción")
	text = strings.ReplaceAll(text, "TransacciÃ³n:", "Transacción:")
	text = strings.ReplaceAll(text, "Valor TransacciÃ³n", "Valor Transacción")
	text = strings.ReplaceAll(text, "Lugar de TransacciÃ³n", "Lugar de Transacción")

	// Normalizar espacios repetidos
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Volver a lines para conservar lógica previa de etiquetas
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
		} else if strings.HasPrefix(line, "Valor Transacción:") || strings.HasPrefix(line, "Valor TransacciÃ³n:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				line = "Valor Transacción: " + value
			}
		} else if strings.HasPrefix(line, "Lugar de Transacción:") || strings.HasPrefix(line, "Lugar de TransacciÃ³n:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				// mantener sin espacio tras ":" si así lo espera la regex; aquí dejamos sin espacio como antes
				line = "Lugar de Transacción:" + value
			}
		}

		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
}

func (d *davivienda) GetMovements() ([]*movementDomain.Movement, error) {
	if d.extract == nil || d.extract.Path == "" || d.password == "" {
		return nil, fmt.Errorf("missing extract, path or password")
	}

	text, err := extractTextFromPDF(d.extract.Path, d.password)
	if err != nil {
		return nil, err
	}

	year, err := parseYear(text)
	if err != nil {
		return nil, err
	}

	matches := movementRegex.FindAllStringSubmatch(text, -1)
	movements := make([]*movementDomain.Movement, 0, len(matches))

	for _, m := range matches {
		mov, err := parseMovement(m, year, d.extract.AccountID, d.extract.MessageID, d.extract.ID)
		if err != nil {
			// Depending on the desired behavior, you might want to log the error and continue
			// or return immediately. For now, we return on the first error.
			return nil, err
		}

		movements = append(movements, mov)
	}

	return movements, nil
}

func parseYear(text string) (int64, error) {
	yearMatch := yearRegex.FindStringSubmatch(text)
	if len(yearMatch) < 2 {
		// Return a default value or an error if the year is not found.
		// For this example, we return 0 and no error, but you might want to handle this differently.
		return 0, nil
	}

	year, err := strconv.ParseInt(yearMatch[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing year: %w", err)
	}

	return year, nil
}

func parseMovement(m []string, year int64, accountID, messageID, extractID string) (*movementDomain.Movement, error) {
	dayAndMonth := strings.Split(m[1], " ")
	date, err := time.Parse("2006-01-02", fmt.Sprintf("%d-%s-%s", year, dayAndMonth[1], dayAndMonth[0]))
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	value, err := strconv.ParseFloat(strings.Replace(m[2], ",", "", -1), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing value: %w", err)
	}

	movementType := movementDomain.Income
	if m[3] == "-" {
		movementType = movementDomain.Expense
	}

	mov := movementDomain.NewMovement(
		accountID,
		institutionID,
		messageID,
		extractID,
		ToValidUTF8(strings.TrimSpace(m[5])),
		value,
		movementDomain.Unknown,
		movementType,
		date,
		movementDomain.ExtractSource,
	)

	return mov, nil
}

func ToValidUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	return string([]rune(s))
}
