package format

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	formatrepo "github.com/m-sharp/edh-tracker/lib/repositories/format"
)

type mockFormatRepo struct {
	GetAllFn    func(ctx context.Context) ([]formatrepo.Model, error)
	GetByIdFn   func(ctx context.Context, id int) (*formatrepo.Model, error)
	GetByNameFn func(ctx context.Context, name string) (*formatrepo.Model, error)
}

func (m *mockFormatRepo) GetAll(ctx context.Context) ([]formatrepo.Model, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	panic("unexpected call to GetAll")
}

func (m *mockFormatRepo) GetById(ctx context.Context, id int) (*formatrepo.Model, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, id)
	}
	panic("unexpected call to GetById")
}

func (m *mockFormatRepo) GetByName(ctx context.Context, name string) (*formatrepo.Model, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, name)
	}
	panic("unexpected call to GetByName")
}

func TestGetAll_Success(t *testing.T) {
	resetCache()
	repo := &mockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			return []formatrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, Name: "commander"},
				{ModelBase: base.ModelBase{ID: 2}, Name: "other"},
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
	assert.True(t, names["commander"])
	assert.True(t, names["other"])
}

func TestGetAll_CachePreventsSecondCall(t *testing.T) {
	resetCache()
	callCount := 0
	repo := &mockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			callCount++
			return []formatrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, Name: "commander"},
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
	repo := &mockFormatRepo{
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
	repo := &mockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			return []formatrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, Name: "commander"},
				{ModelBase: base.ModelBase{ID: 2}, Name: "other"},
			}, nil
		},
	}
	fn := GetByID(repo)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "commander", got.Name)
}

func TestGetByID_NotFound(t *testing.T) {
	resetCache()
	repo := &mockFormatRepo{
		GetAllFn: func(ctx context.Context) ([]formatrepo.Model, error) {
			return []formatrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, Name: "commander"},
			}, nil
		},
	}
	fn := GetByID(repo)
	got, err := fn(context.Background(), 99)
	require.NoError(t, err)
	assert.Nil(t, got)
}
