package pod

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
		{name: "valid", entity: Entity{Name: "My Pod"}, wantErr: false},
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

func TestPlayerPodInputEntityValidate(t *testing.T) {
	tests := []struct {
		name    string
		entity  PlayerPodInputEntity
		wantErr bool
	}{
		{name: "valid", entity: PlayerPodInputEntity{PodID: 1, PlayerID: 2}, wantErr: false},
		{name: "zero pod id", entity: PlayerPodInputEntity{PodID: 0, PlayerID: 2}, wantErr: true},
		{name: "zero player id", entity: PlayerPodInputEntity{PodID: 1, PlayerID: 0}, wantErr: true},
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

func TestJoinInputEntityValidate(t *testing.T) {
	tests := []struct {
		name    string
		entity  JoinInputEntity
		wantErr bool
	}{
		{name: "valid", entity: JoinInputEntity{InviteCode: "abc-123"}, wantErr: false},
		{name: "empty code", entity: JoinInputEntity{InviteCode: ""}, wantErr: true},
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

func TestLeaveInputEntityValidate(t *testing.T) {
	tests := []struct {
		name    string
		entity  LeaveInputEntity
		wantErr bool
	}{
		{name: "valid", entity: LeaveInputEntity{PodID: 1}, wantErr: false},
		{name: "zero pod id", entity: LeaveInputEntity{PodID: 0}, wantErr: true},
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

func TestUpdatePodInputEntityValidate(t *testing.T) {
	tests := []struct {
		name    string
		entity  UpdatePodInputEntity
		wantErr bool
	}{
		{name: "valid", entity: UpdatePodInputEntity{PodID: 1, Name: "My Pod"}, wantErr: false},
		{name: "zero pod id", entity: UpdatePodInputEntity{PodID: 0, Name: "My Pod"}, wantErr: true},
		{name: "empty name", entity: UpdatePodInputEntity{PodID: 1, Name: ""}, wantErr: true},
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
