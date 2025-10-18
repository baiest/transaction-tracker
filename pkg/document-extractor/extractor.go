package documentextractor

import (
	"os/exec"
	"path/filepath"
)

type DaviviendaExtract struct {
	Password string
}

var (
	currentCommand = exec.Command
)

func ExtractTextFromPDF(pathPDF string, password string) (string, error) {
	script := filepath.Join("C:\\Users\\Juan\\Proyectos mios\\transaction-tracker\\pkg\\document-extractor\\dist\\extract.exe")

	cmd := currentCommand(script, pathPDF, password)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
