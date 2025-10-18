package usecase

import (
	"context"
	"errors"
	"testing"

	accountsDomain "transaction-tracker/internal/accounts/domain"
	extractsDomain "transaction-tracker/internal/extracts/domain"

	"github.com/stretchr/testify/require"
)

func TestMockExtractsUsecase_GetByMessageID(t *testing.T) {
	ctx := context.Background()
	mockUsecase := new(MockExtractsUsecase)

	expectedExtract := &extractsDomain.Extract{ID: "ext-123"}
	mockUsecase.
		On("GetByMessageID", ctx, "msg-001").
		Return(expectedExtract, nil)

	extract, err := mockUsecase.GetByMessageID(ctx, "msg-001")

	require.NoError(t, err)
	require.Equal(t, expectedExtract, extract)
	mockUsecase.AssertCalled(t, "GetByMessageID", ctx, "msg-001")
	mockUsecase.AssertExpectations(t)
}

func TestMockExtractsUsecase_GetByMessageID_Error(t *testing.T) {
	ctx := context.Background()
	mockUsecase := new(MockExtractsUsecase)

	mockUsecase.
		On("GetByMessageID", ctx, "msg-err").
		Return(nil, errors.New("not found"))

	extract, err := mockUsecase.GetByMessageID(ctx, "msg-err")

	require.Error(t, err)
	require.Nil(t, extract)
	require.EqualError(t, err, "not found")
	mockUsecase.AssertExpectations(t)
}

func TestMockExtractsUsecase_GetExtractMessages(t *testing.T) {
	ctx := context.Background()
	mockUsecase := new(MockExtractsUsecase)
	account := &accountsDomain.Account{ID: "acc-001"}

	expectedMessages := []string{"msg-1", "msg-2"}
	mockUsecase.
		On("GetExtractMessages", ctx, "MyBank", account).
		Return(expectedMessages, nil)

	msgs, err := mockUsecase.GetExtractMessages(ctx, "MyBank", account)

	require.NoError(t, err)
	require.Equal(t, expectedMessages, msgs)
	mockUsecase.AssertCalled(t, "GetExtractMessages", ctx, "MyBank", account)
	mockUsecase.AssertExpectations(t)
}

func TestMockExtractsUsecase_GetExtractMessages_Error(t *testing.T) {
	ctx := context.Background()
	mockUsecase := new(MockExtractsUsecase)
	account := &accountsDomain.Account{ID: "acc-002"}

	mockUsecase.
		On("GetExtractMessages", ctx, "OtherBank", account).
		Return(nil, errors.New("connection error"))

	msgs, err := mockUsecase.GetExtractMessages(ctx, "OtherBank", account)

	require.Error(t, err)
	require.Nil(t, msgs)
	require.EqualError(t, err, "connection error")
	mockUsecase.AssertExpectations(t)
}

func TestMockExtractsUsecase_Save(t *testing.T) {
	ctx := context.Background()
	mockUsecase := new(MockExtractsUsecase)
	extract := &extractsDomain.Extract{ID: "ext-save"}

	mockUsecase.
		On("Save", ctx, extract).
		Return(nil)

	err := mockUsecase.Save(ctx, extract)

	require.NoError(t, err)
	mockUsecase.AssertCalled(t, "Save", ctx, extract)
	mockUsecase.AssertExpectations(t)
}

func TestMockExtractsUsecase_Save_Error(t *testing.T) {
	ctx := context.Background()
	mockUsecase := new(MockExtractsUsecase)
	extract := &extractsDomain.Extract{ID: "ext-fail"}

	mockUsecase.
		On("Save", ctx, extract).
		Return(errors.New("db error"))

	err := mockUsecase.Save(ctx, extract)

	require.Error(t, err)
	require.EqualError(t, err, "db error")
	mockUsecase.AssertExpectations(t)
}

func TestMockExtractsUsecase_Update(t *testing.T) {
	ctx := context.Background()
	mockUsecase := new(MockExtractsUsecase)
	extract := &extractsDomain.Extract{ID: "ext-update"}

	mockUsecase.
		On("Update", ctx, extract).
		Return(nil)

	err := mockUsecase.Update(ctx, extract)

	require.NoError(t, err)
	mockUsecase.AssertCalled(t, "Update", ctx, extract)
	mockUsecase.AssertExpectations(t)
}

func TestMockExtractsUsecase_Update_Error(t *testing.T) {
	ctx := context.Background()
	mockUsecase := new(MockExtractsUsecase)
	extract := &extractsDomain.Extract{ID: "ext-update-err"}

	mockUsecase.
		On("Update", ctx, extract).
		Return(errors.New("update failed"))

	err := mockUsecase.Update(ctx, extract)

	require.Error(t, err)
	require.EqualError(t, err, "update failed")
	mockUsecase.AssertExpectations(t)
}
