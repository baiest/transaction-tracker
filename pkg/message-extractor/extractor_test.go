package messageextractor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMovementExtractor(t *testing.T) {
	t.Run("returns an error if message type is movement and body is empty", func(t *testing.T) {
		c := require.New(t)

		extractor, err := NewMovementExtractor("davivienda", "", Movement)

		c.Error(err)
		c.Nil(extractor)
		c.Contains(err.Error(), "missing body")
	})

	t.Run("returns a davivienda extractor", func(t *testing.T) {
		c := require.New(t)

		extractor, err := NewMovementExtractor("davivienda", "test body", Movement)

		c.NoError(err)
		c.NotNil(extractor)
		c.IsType(&davivienda{}, extractor)
	})
}
