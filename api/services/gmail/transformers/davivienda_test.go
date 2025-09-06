package transformers

import (
	"testing"
	"time"
	"transaction-tracker/api/services/gmail/models"

	"github.com/stretchr/testify/require"
)

func TestNewDaviviendaMovementFromText(t *testing.T) {
	c := require.New(t)

	text := `Fecha:2025/08/14
Hora:19:24:35
Valor Transacción: $162,000
Clase de Movimiento: Descuento Transferencia a una llave,
Lugar de Transacción:App Transaccional`

	transfomer := NewDaviviendaTransformer(text, models.Movement)

	movement, err := transfomer.Excecute()
	c.NoError(err)

	expectedDate := time.Date(2025, 8, 14, 19, 24, 35, 0, time.Local)
	c.Equal(movement[0].Date, expectedDate)

	c.Equal(162000.00, movement[0].Value)
	c.Equal("others", movement[0].Topic)
	c.Equal("Descuento Transferencia a una llave App Transaccional", movement[0].Detail)
}

func TestNewDaviviendaMovementFromText_InvalidFormat(t *testing.T) {
	c := require.New(t)

	text := `Invalid text without expected labels`
	transfomer := NewDaviviendaTransformer(text, models.Movement)

	_, err := transfomer.Excecute()
	c.ErrorContains(err, "not found labels")
}
