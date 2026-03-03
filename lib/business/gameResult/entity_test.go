package gameResult

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputEntityValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   InputEntity
		wantErr bool
	}{
		{name: "valid", input: InputEntity{DeckID: 1, Place: 1, Kills: 0}, wantErr: false},
		{name: "zero deck id", input: InputEntity{DeckID: 0, Place: 1, Kills: 0}, wantErr: true},
		{name: "place below 1", input: InputEntity{DeckID: 1, Place: 0, Kills: 0}, wantErr: true},
		{name: "negative kills", input: InputEntity{DeckID: 1, Place: 1, Kills: -1}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
