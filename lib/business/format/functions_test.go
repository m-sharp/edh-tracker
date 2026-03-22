package format

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/testHelpers"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	formatrepo "github.com/m-sharp/edh-tracker/lib/repositories/format"
)

func TestGetAll_Success(t *testing.T) {
	resetCache()
	repo := &testHelpers.MockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			return []formatrepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, Name: "commander"},
				{GormModelBase: base.GormModelBase{ID: 2}, Name: "other"},
			}, nil
		},
	}
	fn := GetAll(repo)
	got, err := fn(context.Background())
	require.NoError(t, err)
	assert.Len(t, got, 2)
	names := map[string]bool{}
	for _, e := range got {
		names[e.Name] = true
	}
	assert.True(t, names["Commander"])
	assert.True(t, names["Other"])
}

func TestGetAll_CachePreventsSecondCall(t *testing.T) {
	resetCache()
	callCount := 0
	repo := &testHelpers.MockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			callCount++
			return []formatrepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, Name: "commander"},
			}, nil
		},
	}
	fn := GetAll(repo)
	_, err := fn(context.Background())
	require.NoError(t, err)
	_, err = fn(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestGetAll_RepoError(t *testing.T) {
	resetCache()
	repo := &testHelpers.MockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			return nil, errors.New("db error")
		},
	}
	fn := GetAll(repo)
	_, err := fn(context.Background())
	assert.Error(t, err)
}

func TestGetByID_Found(t *testing.T) {
	resetCache()
	repo := &testHelpers.MockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			return []formatrepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, Name: "commander"},
				{GormModelBase: base.GormModelBase{ID: 2}, Name: "other"},
			}, nil
		},
	}
	fn := GetByID(repo)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Commander", got.Name)
}

func TestGetByID_NotFound(t *testing.T) {
	resetCache()
	repo := &testHelpers.MockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			return []formatrepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, Name: "commander"},
			}, nil
		},
	}
	fn := GetByID(repo)
	got, err := fn(context.Background(), 99)
	require.NoError(t, err)
	assert.Nil(t, got)
}
