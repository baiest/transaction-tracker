package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"transaction-tracker/api/models"
	"transaction-tracker/internal/messages/domain"
	"transaction-tracker/internal/messages/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMessage_Success(t *testing.T) {
	c := require.New(t)

	messageID := "MID123"
	mockUsecase := new(usecase.MockMessageUsecase)
	mockUsecase.On("GetMessageByIDAndAccountID", mock.Anything, messageID, mock.Anything).Return(&domain.Message{ID: messageID}, nil)
	mockUsecase.On("GetMessage", mock.Anything, messageID, mock.Anything).Return(&domain.Message{ID: messageID}, nil)

	testHandler := NewMessageHandler(mockUsecase)

	ginContext, w := setupTestContext(http.MethodGet, "/messages/MID123", nil)
	ginContext.Params = gin.Params{
		{Key: "id", Value: "MID123"},
	}

	testHandler.GetMessageByID(ginContext)

	c.Equal(http.StatusOK, w.Code)

	var response *models.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	c.NoError(err)

	c.Equal(response.ID, messageID)
}

func TestGetMesage_UsecaseError(t *testing.T) {
	c := require.New(t)

	mockUsecase := new(usecase.MockMessageUsecase)
	mockUsecase.On("GetMessageByIDAndAccountID", mock.Anything, "MID123", mock.Anything).Return(nil, errors.New("database connection failed"))
	mockUsecase.On("GetMessage", mock.Anything, "MID123", mock.Anything).Return(nil, errors.New("database connection failed"))

	testHandler := NewMessageHandler(mockUsecase)

	ginContext, w := setupTestContext(http.MethodGet, "/messages/MID123", nil)
	ginContext.Params = gin.Params{
		{Key: "id", Value: "MID123"},
	}

	testHandler.GetMessageByID(ginContext)

	c.Equal(http.StatusInternalServerError, w.Code)
}
