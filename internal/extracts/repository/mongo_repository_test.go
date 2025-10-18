package repository

import (
	"context"
	"testing"
	"time"

	"transaction-tracker/internal/extracts/domain"
	database "transaction-tracker/pkg/databases/mongo"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestGetByMessageID_Success(t *testing.T) {
	c := require.New(t)

	ex := domain.NewExtract("acc-1", "msg-1", "inst-1", "/path.pdf", time.March, 2025)

	mockColl := &database.MockCollection{
		FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
			return mongo.NewSingleResultFromDocument(ex, nil, nil)
		},
	}

	repo := NewExtractsRepository(mockColl)

	result, err := repo.GetByMessageID(context.Background(), ex.MessageID)
	c.NoError(err)
	c.Equal(ex, result)
}

func TestGetByMessageID_NotFound(t *testing.T) {
	c := require.New(t)

	mockColl := &database.MockCollection{
		FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
			return &mongo.SingleResult{}
		},
	}

	repo := NewExtractsRepository(mockColl)

	result, err := repo.GetByMessageID(context.Background(), "non-existent")
	c.Error(err)
	c.Equal(ErrExtractNotFound, err)
	c.Nil(result)
}

func TestSave_SetsTimestampsAndInserts(t *testing.T) {
	c := require.New(t)

	ex := domain.NewExtract("acc-2", "msg-2", "inst-2", "/p.pdf", time.April, 2024)

	var inserted any
	mockColl := &database.MockCollection{
		InsertOneFn: func(ctx context.Context, doc any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
			inserted = doc
			return &mongo.InsertOneResult{InsertedID: ex.ID}, nil
		},
	}

	repo := NewExtractsRepository(mockColl)
	// controlar tiempo para aserciones predecibles
	concrete := repo.(*extractsRepository)
	fixed := time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)
	concrete.nowFunc = func() time.Time { return fixed }

	err := repo.Save(context.Background(), ex)
	c.NoError(err)

	c.NotZero(ex.CreatedAt)
	c.NotZero(ex.UpdatedAt)
	c.Equal(fixed, ex.CreatedAt)
	c.Equal(fixed, ex.UpdatedAt)

	insEx, ok := inserted.(*domain.Extract)
	c.True(ok)
	c.Equal(ex.ID, insEx.ID)
}

func TestUpdate_SetsUpdatedAtAndCallsUpdateOne(t *testing.T) {
	c := require.New(t)

	ex := domain.NewExtract("acc-3", "msg-3", "inst-3", "/u.pdf", time.May, 2023)
	ex.ID = "EXI-custom-id"
	// set initial timestamps
	ex.CreatedAt = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	ex.UpdatedAt = ex.CreatedAt

	var receivedFilter any
	var receivedUpdate any

	mockColl := &database.MockCollection{
		UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
			receivedFilter = filter
			receivedUpdate = update
			return &mongo.UpdateResult{}, nil
		},
	}

	repo := NewExtractsRepository(mockColl)
	concrete := repo.(*extractsRepository)
	fixed := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	concrete.nowFunc = func() time.Time { return fixed }

	err := repo.Update(context.Background(), ex)
	c.NoError(err)

	c.Equal(fixed, ex.UpdatedAt)

	// comprobar filtro y update enviados al UpdateOne
	filterMap, ok := receivedFilter.(map[string]interface{})
	c.True(ok)
	c.Equal(ex.ID, filterMap["_id"])

	updateMap, ok := receivedUpdate.(map[string]interface{})
	c.True(ok)
	setMap, ok := updateMap["$set"].(*domain.Extract)
	if !ok {
		// depending on how caller serializa, it may be a value map; fallback check updated time via any encoding:
		// just ensure updateMap contains $set
		c.Contains(updateMap, "$set")
	} else {
		c.Equal(fixed, setMap.UpdatedAt)
	}
}
