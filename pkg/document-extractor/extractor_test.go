package documentextractor

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractTextFromPDF(t *testing.T) {
	c := require.New(t)

	prevRunCommand := runCommand
	defer func() { runCommand = prevRunCommand }()

	runCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("cmd", "/C", "echo", "Texto simulado desde PDF")
	}

	out, err := ExtractTextFromPDF("dummy.pdf", "1234")
	c.NoError(err)
	c.Contains(out, "Texto simulado")
}
