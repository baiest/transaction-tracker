package documentextractor

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

//go:embed extract.py
var extractorPy []byte

var (
	runCommand = exec.Command

	scriptPath string
	initOnce   sync.Once
)

// prepareScript crea el script solo una vez y lo reutiliza
func prepareScript() error {
	var err error
	initOnce.Do(func() {
		tmpDir := os.TempDir()
		scriptPath = filepath.Join(tmpDir, "extractor_embedded.py")
		err = os.WriteFile(scriptPath, extractorPy, 0700)
	})

	return err
}

func ExtractTextFromPDF(pathPDF string, password string) (string, error) {
	// Asegurar que el script existe (solo se crea la primera vez)
	if err := prepareScript(); err != nil {
		return "", err
	}

	// Detectar python o python3
	pythonCmd := "python"
	if _, err := exec.LookPath("python3"); err == nil {
		pythonCmd = "python3"
	}

	// Ejecutar el script con los argumentos
	cmd := runCommand(pythonCmd, scriptPath, pathPDF, password)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
