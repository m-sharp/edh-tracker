package pod

import (
	"fmt"
	"time"

	repo "github.com/m-sharp/edh-tracker/lib/repositories/pod"
)

type Entity struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e Entity) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(e.Name) > 255 {
		return fmt.Errorf("name must be 255 characters or fewer")
	}
	return nil
}

// PlayerPodInputEntity is the input for adding a player to a pod.
type PlayerPodInputEntity struct {
	PodID    int `json:"pod_id"`
	PlayerID int `json:"player_id"`
}

func (p PlayerPodInputEntity) Validate() error {
	if p.PodID == 0 {
		return fmt.Errorf("pod_id is required")
	}
	if p.PlayerID == 0 {
		return fmt.Errorf("player_id is required")
	}
	return nil
}

// PlayerWithRole is a pod member with their role.
type PlayerWithRole struct {
	PlayerID int    `json:"player_id"`
	Role     string `json:"role"`
}

// InviteEntity is the response for a generated invite code.
type InviteEntity struct {
	InviteCode string `json:"invite_code"`
}

// JoinInputEntity is the input for joining a pod by invite code.
type JoinInputEntity struct {
	InviteCode string `json:"invite_code"`
}

func (j JoinInputEntity) Validate() error {
	if j.InviteCode == "" {
		return fmt.Errorf("invite_code is required")
	}
	return nil
}

// LeaveInputEntity is the input for leaving a pod.
type LeaveInputEntity struct {
	PodID int `json:"pod_id"`
}

func (l LeaveInputEntity) Validate() error {
	if l.PodID == 0 {
		return fmt.Errorf("pod_id is required")
	}
	return nil
}

// UpdatePodInputEntity is the input for updating a pod name.
type UpdatePodInputEntity struct {
	PodID int    `json:"pod_id"`
	Name  string `json:"name"`
}

func (u UpdatePodInputEntity) Validate() error {
	if u.PodID == 0 {
		return fmt.Errorf("pod_id is required")
	}
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(u.Name) > 255 {
		return fmt.Errorf("name must be 255 characters or fewer")
	}
	return nil
}

func ToEntity(m repo.Model) Entity {
	return Entity{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
