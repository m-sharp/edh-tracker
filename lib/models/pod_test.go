package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPod_Validate(t *testing.T) {
	tests := []struct {
		name    string
		pod     Pod
		wantErr bool
	}{
		{
			name:    "valid",
			pod:     Pod{Name: "My Pod"},
			wantErr: false,
		},
		{
			name:    "missing Name",
			pod:     Pod{Name: ""},
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(tt *testing.T) {
			err := testCase.pod.Validate()
			if testCase.wantErr {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
		})
	}
}

func TestPlayerPod_Validate(t *testing.T) {
	tests := []struct {
		name      string
		playerPod PlayerPod
		wantErr   bool
	}{
		{
			name:      "valid",
			playerPod: PlayerPod{PodID: 1, PlayerID: 2},
			wantErr:   false,
		},
		{
			name:      "missing PodID",
			playerPod: PlayerPod{PodID: 0, PlayerID: 2},
			wantErr:   true,
		},
		{
			name:      "missing PlayerID",
			playerPod: PlayerPod{PodID: 1, PlayerID: 0},
			wantErr:   true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(tt *testing.T) {
			err := testCase.playerPod.Validate()
			if testCase.wantErr {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
		})
	}
}
