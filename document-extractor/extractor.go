package documentextractor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	currentTextExtractor = extractTextFromPDF
)

type Movement struct {
	Date       time.Time
	Value      float64
	IsNegative bool
	Type       string
	Detail     string
}

type DaviviendaExtract struct {
	Password string
}

func extractTextFromPDF(pathPDF string, password string) string {
	exePath, _ := os.Getwd()
	dir := filepath.Dir(exePath)

	script := filepath.Join(dir, "document-extractor\\dist\\extract.exe")
	pdf := filepath.Join(dir, pathPDF)

	cmd := exec.Command(script, pdf, password)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return string(out)
}

func (e *DaviviendaExtract) GetMovements(pathPDF string) []*Movement {
	text := currentTextExtractor(pathPDF, e.Password)

	regex := regexp.MustCompile(`(?m)^(\d{2}\s+\d{2})\s+\$\s*([\d,]+\.\d{2})([+-])\s+(\d{4})\s+(.+)$`)
	matches := regex.FindAllStringSubmatch(text, -1)

	movements := []*Movement{}

	for _, m := range matches {
		dayAndMonth := strings.Split(m[1], " ")

		date, err := time.Parse("2006-01-02", fmt.Sprintf("2021-%s-%s", dayAndMonth[1], dayAndMonth[0]))
		if err != nil {
			fmt.Println(err)
			continue
		}

		value, err := strconv.ParseFloat(strings.Replace(m[2], ",", "", -1), 64)
		if err != nil {
			fmt.Println(err)

			continue
		}

		mov := &Movement{
			Date:       date,
			Value:      value,
			IsNegative: m[3] == "-",
			Type:       "unknown",
			Detail:     ToValidUTF8(strings.TrimSpace(m[5])),
		}

		movements = append(movements, mov)
	}

	return movements
}

func ToValidUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	// Reemplaza bytes inválidos por el caracter de sustitución
	return string([]rune(s))
}
