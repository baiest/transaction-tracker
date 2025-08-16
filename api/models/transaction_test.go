package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewDaviviendaMovementFromText(t *testing.T) {
	c := require.New(t)

	text := `Fecha:2025/08/14
Hora:19:24:35
Valor Transacción: $162,000
Clase de Movimiento: Descuento Transferencia a una llave,
Lugar de Transacción:App Transaccional`

	movement, err := NewDaviviendaMovementFromText(text)
	c.NoError(err)
	c.NotNil(movement)

	expectedDate := time.Date(2025, 8, 14, 0, 0, 0, 0, time.Local)
	c.Equal(movement.Date, expectedDate)

	c.Equal(19, movement.Hour.Hour())
	c.Equal(24, movement.Hour.Minute())
	c.Equal(35, movement.Hour.Second())

	c.Equal(162000.00, movement.Value)
	c.Equal("Descuento Transferencia a una llave", movement.Type)
	c.Equal("App Transaccional", movement.Place)
}

func TestNewDaviviendaMovementFromText_InvalidFormat(t *testing.T) {
	c := require.New(t)

	text := `Invalid text without expected labels`
	movement, err := NewDaviviendaMovementFromText(text)
	c.Error(err)
	c.Nil(movement)
}

func TestTransactionRequest_MarshalJSON(t *testing.T) {
	c := require.New(t)

	date := time.Date(2025, 8, 14, 19, 24, 35, 0, time.FixedZone("EST", -5*60*60))
	tr := &TransactionRequest{
		Date:  date,
		Type:  "Descuento Transferencia a una llave",
		Value: 162000,
	}

	jsonData, err := json.Marshal(tr)
	c.NoError(err)
	c.NotNil(jsonData)

	jsonStr := string(jsonData)
	c.Contains(jsonStr, `"date":"2025-08-15T00:24:35Z"`) // UTC time (EST + 5 hours)
	c.Contains(jsonStr, `"type":"Descuento Transferencia a una llave"`)
	c.Contains(jsonStr, `"value":162000`)
}
