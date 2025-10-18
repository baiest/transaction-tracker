package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/api/gmail/v1"

	accountsDomain "transaction-tracker/internal/accounts/domain"
	messageextractor "transaction-tracker/pkg/message-extractor"
)

func TestIsMessageFiltered(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name       string
		headers    []*gmail.MessagePartHeader
		wantType   messageextractor.MessageType
		wantFilter bool
	}{
		{
			name: "davivienda movimiento",
			headers: []*gmail.MessagePartHeader{
				{Name: "From", Value: "banco_davivienda@davivienda.com"},
				{Name: "Subject", Value: "Davivienda"},
			},
			wantType:   messageextractor.Movement,
			wantFilter: true,
		},
		{
			name: "extracto",
			headers: []*gmail.MessagePartHeader{
				{Name: "From", Value: "bancodavivienda@davivienda.com"},
				{Name: "Subject", Value: "Extractos Septiembre"},
			},
			wantType:   messageextractor.Extract,
			wantFilter: true,
		},
		{
			name: "unknown sender",
			headers: []*gmail.MessagePartHeader{
				{Name: "From", Value: "otro@banco.com"},
				{Name: "Subject", Value: "Notificación"},
			},
			wantType:   messageextractor.Unknown,
			wantFilter: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &gmail.Message{Payload: &gmail.MessagePart{Headers: tt.headers}}
			gotType, gotFilter := isMessageFiltered(msg)
			c.Equal(tt.wantType, gotType)
			c.Equal(tt.wantFilter, gotFilter)
		})
	}
}

func TestProcess_MissingExternalID(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	u := &messageUsecase{}

	account := &accountsDomain.Account{ID: "acc1"}

	msg, err := u.Process(ctx, "notif1", "", account)
	c.Nil(msg)
	c.ErrorIs(err, ErrMissingExternalID)
}
