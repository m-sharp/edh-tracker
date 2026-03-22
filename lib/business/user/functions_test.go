package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/testHelpers"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	userrepo "github.com/m-sharp/edh-tracker/lib/repositories/user"
)

const (
	testEmail       = "test@example.com"
	testProvider    = "google"
	testSubject     = "sub-123"
	testDisplayName = "Test User"
	testAvatarURL   = "https://example.com/avatar.png"
)

func TestGetByEmail_Found(t *testing.T) {
	repo := &testHelpers.MockUserRepo{
		GetByEmailFn: func(_ context.Context, e string) (*userrepo.Model, error) {
			return &userrepo.Model{GormModelBase: base.GormModelBase{ID: 7}, PlayerID: 3, Email: &e}, nil
		},
	}

	fn := GetByEmail(repo)
	got, err := fn(context.Background(), testEmail)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 7, got.ID)
	assert.Equal(t, 3, got.PlayerID)
}

func TestGetByEmail_NotFound(t *testing.T) {
	repo := &testHelpers.MockUserRepo{
		GetByEmailFn: func(_ context.Context, _ string) (*userrepo.Model, error) {
			return nil, nil
		},
	}

	fn := GetByEmail(repo)
	got, err := fn(context.Background(), "nobody@nowhere.example")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByEmail_Error(t *testing.T) {
	repo := &testHelpers.MockUserRepo{
		GetByEmailFn: func(_ context.Context, _ string) (*userrepo.Model, error) {
			return nil, errors.New("db error")
		},
	}

	fn := GetByEmail(repo)
	_, err := fn(context.Background(), testEmail)
	assert.Error(t, err)
}

func TestLinkOAuth_Success(t *testing.T) {
	repo := &testHelpers.MockUserRepo{
		UpdateOAuthFn: func(_ context.Context, _ int, _, _, _, _, _ string) error {
			return nil
		},
		GetByIDFn: func(_ context.Context, id int) (*userrepo.Model, error) {
			sub := testSubject
			email := testEmail
			return &userrepo.Model{
				GormModelBase: base.GormModelBase{ID: id},
				PlayerID:      5,
				OAuthSubject:  &sub,
				Email:         &email,
			}, nil
		},
	}

	fn := LinkOAuth(repo)
	got, err := fn(context.Background(), 10, testProvider, testSubject, testEmail, testDisplayName, testAvatarURL)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 10, got.ID)
	assert.Equal(t, 5, got.PlayerID)
}

func TestLinkOAuth_UpdateError(t *testing.T) {
	repo := &testHelpers.MockUserRepo{
		UpdateOAuthFn: func(_ context.Context, _ int, _, _, _, _, _ string) error {
			return errors.New("update error")
		},
	}

	fn := LinkOAuth(repo)
	_, err := fn(context.Background(), 10, testProvider, testSubject, testEmail, testDisplayName, testAvatarURL)
	assert.Error(t, err)
}

func TestLinkOAuth_FetchError(t *testing.T) {
	repo := &testHelpers.MockUserRepo{
		UpdateOAuthFn: func(_ context.Context, _ int, _, _, _, _, _ string) error {
			return nil
		},
		GetByIDFn: func(_ context.Context, _ int) (*userrepo.Model, error) {
			return nil, errors.New("fetch error")
		},
	}

	fn := LinkOAuth(repo)
	_, err := fn(context.Background(), 10, testProvider, testSubject, testEmail, testDisplayName, testAvatarURL)
	assert.Error(t, err)
}
