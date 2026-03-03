package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreate(t *testing.T) {
	tests := []struct {
		name     string
		playerID int
		deckName string
		formatID int
		wantErr  bool
	}{
		{name: "valid", playerID: 1, deckName: "Test Deck", formatID: 1, wantErr: false},
		{name: "zero player id", playerID: 0, deckName: "Test Deck", formatID: 1, wantErr: true},
		{name: "empty name", playerID: 1, deckName: "", formatID: 1, wantErr: true},
		{name: "zero format id", playerID: 1, deckName: "Test Deck", formatID: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreate(tt.playerID, tt.deckName, tt.formatID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
