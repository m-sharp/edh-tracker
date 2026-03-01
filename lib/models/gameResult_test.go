package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameResult_Validate(t *testing.T) {
	tests := []struct {
		name    string
		result  GameResult
		wantErr bool
	}{
		{
			name:    "valid",
			result:  GameResult{DeckId: 1, Place: 2, Kills: 0},
			wantErr: false,
		},
		{
			name:    "missing DeckId",
			result:  GameResult{DeckId: 0, Place: 1, Kills: 0},
			wantErr: true,
		},
		{
			name:    "missing Place (zero)",
			result:  GameResult{DeckId: 1, Place: 0, Kills: 0},
			wantErr: true,
		},
		{
			name:    "negative Place",
			result:  GameResult{DeckId: 1, Place: -1, Kills: 0},
			wantErr: true,
		},
		{
			name:    "negative Kills",
			result:  GameResult{DeckId: 1, Place: 1, Kills: -1},
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(tt *testing.T) {
			err := testCase.result.Validate()
			if testCase.wantErr {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
		})
	}
}
