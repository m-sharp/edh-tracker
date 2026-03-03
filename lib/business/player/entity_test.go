package player

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityValidate(t *testing.T) {
	tests := []struct {
		name    string
		entity  Entity
		wantErr bool
	}{
		{name: "valid", entity: Entity{Name: "Alice"}, wantErr: false},
		{name: "empty name", entity: Entity{Name: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entity.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
