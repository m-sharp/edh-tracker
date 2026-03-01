package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		player  Player
		wantErr bool
	}{
		{
			name:    "valid",
			player:  Player{Name: "Alice"},
			wantErr: false,
		},
		{
			name:    "missing Name",
			player:  Player{Name: ""},
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(tt *testing.T) {
			err := testCase.player.Validate()
			if testCase.wantErr {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
		})
	}
}
