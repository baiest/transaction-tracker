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
	Topic      string
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
	fmt.Println("Extracted text length:", len(text))

	reYear := regexp.MustCompile(`INFORME DEL MES:.*?/(\d{4})`)
	yearMatch := reYear.FindStringSubmatch(text)
	year := int64(0)
	var err error

	if len(yearMatch) > 1 {
		year, err = strconv.ParseInt(yearMatch[1], 10, 64)
		if err != nil {
			year = 0
		}
	}

	regex := regexp.MustCompile(`(?m)^(\d{2}\s+\d{2})\s+\$\s*([\d,]+\.\d{2})([+-])\s+(\d{4})\s+(.+)$`)
	matches := regex.FindAllStringSubmatch(text, -1)

	movements := []*Movement{}

	for _, m := range matches {
		fmt.Println("Match:", m)
		dayAndMonth := strings.Split(m[1], " ")

		date, err := time.Parse("2006-01-02", fmt.Sprintf("%d-%s-%s", year, dayAndMonth[1], dayAndMonth[0]))
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
			Topic:      "unknown",
			Detail:     ToValidUTF8(strings.TrimSpace(m[5])),
		}

		fmt.Println("Movement:", mov)

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
