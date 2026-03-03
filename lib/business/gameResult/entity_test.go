package gameResult

import (
	"testing"

	"github.com/m-sharp/edh-tracker/lib/business/stats"
	"github.com/stretchr/testify/assert"
)

func TestGetPointsForPlace(t *testing.T) {
	tests := []struct {
		kills, place, want int
	}{
		{kills: 2, place: 1, want: 5}, // 2 kills + 3 bonus
		{kills: 0, place: 2, want: 2}, // 0 kills + 2 bonus
		{kills: 1, place: 3, want: 2}, // 1 kill + 1 bonus
		{kills: 3, place: 4, want: 3}, // 3 kills + 0 bonus
		{kills: 0, place: 0, want: 0}, // 0 kills + 0 bonus
	}
	for _, tt := range tests {
		got := stats.GetPointsForPlace(tt.kills, tt.place)
		assert.Equal(t, tt.want, got)
	}
}

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
