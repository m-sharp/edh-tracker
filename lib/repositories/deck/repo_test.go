package deck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/deck"
	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetAll_PopulatesAssociations(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	dcRepo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()

	formatID := testHelpers.GetCommanderFormatID(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)
	cmdID := testHelpers.CreateTestCommander(t, db)
	partnerID := testHelpers.CreateTestCommander(t, db)

	deckID, err := deckRepo.Add(ctx, playerID, "Krenko Goblins", formatID)
	require.NoError(t, err)

	_, err = dcRepo.Add(ctx, deckID, cmdID, &partnerID)
	require.NoError(t, err)

	got, err := deckRepo.GetAll(ctx)
	require.NoError(t, err)

	var found *deck.Model
	for i := range got {
		if got[i].ID == deckID {
			found = &got[i]
			break
		}
	}
	require.NotNil(t, found, "deck should be in GetAll results")

	assert.NotEmpty(t, found.Player.Name, "Player.Name should be populated")
	assert.NotEmpty(t, found.Format.Name, "Format.Name should be populated")
	require.NotNil(t, found.Commander, "Commander should be populated")
	assert.Equal(t, cmdID, found.Commander.CommanderID)
	assert.NotEmpty(t, found.Commander.Commander.Name, "Commander.Commander.Name should be populated")
	require.NotNil(t, found.Commander.PartnerCommanderID)
	assert.Equal(t, partnerID, *found.Commander.PartnerCommanderID)
	require.NotNil(t, found.Commander.PartnerCommander)
	assert.NotEmpty(t, found.Commander.PartnerCommander.Name, "PartnerCommander.Name should be populated")
}

func TestGetAll_NoCommander(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	deckID, err := deckRepo.Add(ctx, playerID, "No Commander Deck", formatID)
	require.NoError(t, err)

	got, err := deckRepo.GetAll(ctx)
	require.NoError(t, err)

	var found *deck.Model
	for i := range got {
		if got[i].ID == deckID {
			found = &got[i]
			break
		}
	}
	require.NotNil(t, found)
	assert.Nil(t, found.Commander, "Commander should be nil when no deck_commander row exists")
	assert.NotEmpty(t, found.Player.Name)
	assert.NotEmpty(t, found.Format.Name)
}

func TestGetAll_ExcludesRetired(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	activeID, err := deckRepo.Add(ctx, playerID, "Active", formatID)
	require.NoError(t, err)
	retiredID, err := deckRepo.Add(ctx, playerID, "Retired", formatID)
	require.NoError(t, err)
	require.NoError(t, deckRepo.Retire(ctx, retiredID))

	got, err := deckRepo.GetAll(ctx)
	require.NoError(t, err)

	var ids []int
	for _, d := range got {
		ids = append(ids, d.ID)
	}
	assert.Contains(t, ids, activeID)
	assert.NotContains(t, ids, retiredID)
}

func TestGetAllForPlayer_PopulatesAssociations(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	dcRepo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	cmdID := testHelpers.CreateTestCommander(t, db)

	deckID, err := deckRepo.Add(ctx, playerID, "Test Deck", formatID)
	require.NoError(t, err)
	_, err = dcRepo.Add(ctx, deckID, cmdID, nil)
	require.NoError(t, err)

	got, err := deckRepo.GetAllForPlayer(ctx, playerID)
	require.NoError(t, err)
	require.Len(t, got, 1)

	d := got[0]
	assert.Equal(t, deckID, d.ID)
	assert.NotEmpty(t, d.Player.Name)
	assert.NotEmpty(t, d.Format.Name)
	require.NotNil(t, d.Commander)
	assert.Equal(t, cmdID, d.Commander.CommanderID)
	assert.NotEmpty(t, d.Commander.Commander.Name)
	assert.Nil(t, d.Commander.PartnerCommander)
}

func TestGetAllByPlayerIDs_PopulatesAssociations(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)
	p3 := testHelpers.CreateTestPlayer(t, db)

	id1, err := deckRepo.Add(ctx, p1, "P1 Deck", formatID)
	require.NoError(t, err)
	id2, err := deckRepo.Add(ctx, p2, "P2 Deck", formatID)
	require.NoError(t, err)
	id3, err := deckRepo.Add(ctx, p3, "P3 Deck", formatID)
	require.NoError(t, err)

	got, err := deckRepo.GetAllByPlayerIDs(ctx, []int{p1, p2})
	require.NoError(t, err)

	var ids []int
	for _, d := range got {
		ids = append(ids, d.ID)
		assert.NotEmpty(t, d.Player.Name, "Player.Name should be populated")
		assert.NotEmpty(t, d.Format.Name, "Format.Name should be populated")
	}
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2)
	assert.NotContains(t, ids, id3)
}

func TestGetAllByPlayerIDs_Empty(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	got, err := deckRepo.GetAllByPlayerIDs(context.Background(), []int{})
	require.NoError(t, err)
	assert.Len(t, got, 0)
}

func TestGetById_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "Krenko Goblins", formatID)
	require.NoError(t, err)

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, playerID, got.PlayerID)
	assert.Equal(t, "Krenko Goblins", got.Name)
	assert.Equal(t, formatID, got.FormatID)
	assert.False(t, got.Retired)
}

func TestGetById_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	got, err := repo.GetById(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByIDHydrated_PopulatesAssociations(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	dcRepo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	cmdID := testHelpers.CreateTestCommander(t, db)

	deckID, err := deckRepo.Add(ctx, playerID, "Krenko Goblins", formatID)
	require.NoError(t, err)
	_, err = dcRepo.Add(ctx, deckID, cmdID, nil)
	require.NoError(t, err)

	got, err := deckRepo.GetByIDHydrated(ctx, deckID)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, deckID, got.ID)
	assert.NotEmpty(t, got.Player.Name)
	assert.NotEmpty(t, got.Format.Name)
	require.NotNil(t, got.Commander)
	assert.Equal(t, cmdID, got.Commander.CommanderID)
	assert.NotEmpty(t, got.Commander.Commander.Name)
	assert.Nil(t, got.Commander.PartnerCommander)
}

func TestGetByIDHydrated_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	got, err := deckRepo.GetByIDHydrated(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByIDHydrated_WithPartner(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	deckRepo := testHelpers.NewDeckRepo(db)
	dcRepo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	cmdID := testHelpers.CreateTestCommander(t, db)
	partnerID := testHelpers.CreateTestCommander(t, db)

	deckID, err := deckRepo.Add(ctx, playerID, "Partner Deck", formatID)
	require.NoError(t, err)
	_, err = dcRepo.Add(ctx, deckID, cmdID, &partnerID)
	require.NoError(t, err)

	got, err := deckRepo.GetByIDHydrated(ctx, deckID)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.NotNil(t, got.Commander)
	assert.NotEmpty(t, got.Commander.Commander.Name)
	require.NotNil(t, got.Commander.PartnerCommanderID)
	assert.Equal(t, partnerID, *got.Commander.PartnerCommanderID)
	require.NotNil(t, got.Commander.PartnerCommander)
	assert.NotEmpty(t, got.Commander.PartnerCommander.Name)
}

func TestAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	formatID := testHelpers.GetCommanderFormatID(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(context.Background(), playerID, "New Deck", formatID)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestBulkAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)

	decks := []deck.Model{
		{PlayerID: p1, Name: "Bulk Deck A", FormatID: formatID},
		{PlayerID: p1, Name: "Bulk Deck B", FormatID: formatID},
		{PlayerID: p2, Name: "Bulk Deck C", FormatID: formatID},
	}
	got, err := repo.BulkAdd(ctx, decks)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	for _, d := range got {
		assert.Greater(t, d.ID, 0)
	}
}

func TestUpdate_PartialFields(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "Original Name", formatID)
	require.NoError(t, err)

	newName := "Updated Name"
	require.NoError(t, repo.Update(ctx, id, deck.UpdateFields{Name: &newName}))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, formatID, got.FormatID) // unchanged
}

func TestUpdate_MultipleFields(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "Original Name", formatID)
	require.NoError(t, err)

	newName := "Updated Name"
	require.NoError(t, repo.Update(ctx, id, deck.UpdateFields{
		Name:    &newName,
		Retired: new(true),
	}))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, newName, got.Name)
	assert.Equal(t, formatID, got.FormatID) // unchanged
	assert.True(t, got.Retired)
}

func TestUpdate_NoFields(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "No Change Deck", formatID)
	require.NoError(t, err)

	require.NoError(t, repo.Update(ctx, id, deck.UpdateFields{}))
}

func TestUpdate_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	name := "Ghost"
	err := repo.Update(context.Background(), 999999, deck.UpdateFields{Name: &name})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected rows")
}

func TestRetire(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "To Retire", formatID)
	require.NoError(t, err)

	require.NoError(t, repo.Retire(ctx, id))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.Retired)
}

func TestRetire_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	err := repo.Retire(context.Background(), 999999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected rows")
}

func TestSoftDelete(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "To Delete", formatID)
	require.NoError(t, err)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}
