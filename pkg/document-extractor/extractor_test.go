package documentextractor

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)



func TestDaviviendaExtractor(t *testing.T) {
	c := require.New(t)

	prevCommand := currentCommand
	defer func() { currentCommand = prevCommand }()

	currentCommand = func(name string, arg ...string) *exec.Cmd {
		return &exec.Cmd{
			Path:   name,
			Args:   append([]string{name}, arg...),
			Stdout: nil,
		}
	}

	_, err := ExtractTextFromPDF("path", "")
	c.NoError(err)
}
