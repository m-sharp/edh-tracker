package models

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeckRow_ToDeck(t *testing.T) {
	base := deckRow{
		PlayerID:   1,
		PlayerName: "Alice",
		Name:       "Deck Alpha",
		FormatID:   1,
		FormatName: "commander",
		Retired:    false,
	}

	t.Run("no commander", func(tt *testing.T) {
		deck := base.toDeck()
		assert.Nil(tt, deck.Commanders)
	})

	t.Run("solo commander", func(tt *testing.T) {
		row := base
		row.CommanderID = sql.NullInt64{Valid: true, Int64: 42}
		row.CommanderName = sql.NullString{Valid: true, String: "Atraxa"}
		deck := row.toDeck()
		assert.NotNil(tt, deck.Commanders)
		assert.Equal(tt, 42, deck.Commanders.CommanderID)
		assert.Equal(tt, "Atraxa", deck.Commanders.CommanderName)
		assert.Nil(tt, deck.Commanders.PartnerCommanderID)
		assert.Nil(tt, deck.Commanders.PartnerCommanderName)
	})

	t.Run("partner commanders", func(tt *testing.T) {
		row := base
		row.CommanderID = sql.NullInt64{Valid: true, Int64: 42}
		row.CommanderName = sql.NullString{Valid: true, String: "Rograkh"}
		row.PartnerCommanderID = sql.NullInt64{Valid: true, Int64: 99}
		row.PartnerCommanderName = sql.NullString{Valid: true, String: "Silas Renn"}
		deck := row.toDeck()
		assert.NotNil(tt, deck.Commanders)
		assert.Equal(tt, 42, deck.Commanders.CommanderID)
		assert.NotNil(tt, deck.Commanders.PartnerCommanderID)
		assert.Equal(tt, 99, *deck.Commanders.PartnerCommanderID)
		assert.NotNil(tt, deck.Commanders.PartnerCommanderName)
		assert.Equal(tt, "Silas Renn", *deck.Commanders.PartnerCommanderName)
	})
}

func TestDeck_Validate(t *testing.T) {
	tests := []struct {
		name    string
		deck    Deck
		wantErr bool
	}{
		{
			name:    "valid",
			deck:    Deck{PlayerID: 1, Name: "My Deck", FormatID: 1},
			wantErr: false,
		},
		{
			name:    "missing PlayerID",
			deck:    Deck{PlayerID: 0, Name: "My Deck", FormatID: 1},
			wantErr: true,
		},
		{
			name:    "missing Name",
			deck:    Deck{PlayerID: 1, Name: "", FormatID: 1},
			wantErr: true,
		},
		{
			name:    "missing FormatID",
			deck:    Deck{PlayerID: 1, Name: "My Deck", FormatID: 0},
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(tt *testing.T) {
			err := testCase.deck.Validate()
			if testCase.wantErr {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
		})
	}
}
