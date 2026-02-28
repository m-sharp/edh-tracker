package models

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetAllPods        = `SELECT id, name, created_at, updated_at, deleted_at FROM pod WHERE deleted_at IS NULL;`
	GetPodByID        = `SELECT id, name, created_at, updated_at, deleted_at FROM pod WHERE id = ? AND deleted_at IS NULL;`
	GetPodsByPlayerID = `SELECT pod.id, pod.name, pod.created_at, pod.updated_at, pod.deleted_at
						   FROM (pod INNER JOIN player_pod ON pod.id = player_pod.pod_id)
						  WHERE player_pod.player_id = ?
						    AND pod.deleted_at IS NULL
						    AND player_pod.deleted_at IS NULL;`
	GetPodIDsByPlayerID = `SELECT pod_id FROM player_pod WHERE player_id = ? AND deleted_at IS NULL;`

	InsertPod       = `INSERT INTO pod (name) VALUES (?);`
	InsertPlayerPod = `INSERT INTO player_pod (pod_id, player_id) VALUES (?, ?);`

	podValidationErr       = "invalid Pod: %s"
	playerPodValidationErr = "invalid PlayerPod: %s"
)

type Pod struct {
	Model
	Name string `json:"name" db:"name"`
}

type PlayerPod struct {
	Model
	PodID    int `json:"pod_id"    db:"pod_id"`
	PlayerID int `json:"player_id" db:"player_id"`
}

func (p *Pod) Validate() error {
	if p.Name == "" {
		return fmt.Errorf(podValidationErr, "missing Name")
	}
	return nil
}

func (pp *PlayerPod) Validate() error {
	if pp.PodID == 0 {
		return fmt.Errorf(playerPodValidationErr, "missing PodID")
	}
	if pp.PlayerID == 0 {
		return fmt.Errorf(playerPodValidationErr, "missing PlayerID")
	}
	return nil
}

type PodProvider struct {
	client *lib.DBClient
}

func NewPodProvider(client *lib.DBClient) *PodProvider {
	return &PodProvider{client: client}
}

func (p *PodProvider) GetAll(ctx context.Context) ([]Pod, error) {
	var pods []Pod
	if err := p.client.Db.SelectContext(ctx, &pods, GetAllPods); err != nil {
		return nil, fmt.Errorf("failed to get Pod records: %w", err)
	}
	if pods == nil {
		return []Pod{}, nil
	}
	return pods, nil
}

func (p *PodProvider) GetByID(ctx context.Context, podID int) (*Pod, error) {
	var pods []Pod
	if err := p.client.Db.SelectContext(ctx, &pods, GetPodByID, podID); err != nil {
		return nil, fmt.Errorf("failed to get Pod record for id %d: %w", podID, err)
	}

	if len(pods) == 0 || len(pods) > 1 {
		return nil, fmt.Errorf(
			"unexpected number of pods returned for ID %d: got %d, expected 1",
			podID, len(pods),
		)
	}

	return &pods[0], nil
}

func (p *PodProvider) GetByPlayerID(ctx context.Context, playerID int) ([]Pod, error) {
	var pods []Pod
	if err := p.client.Db.SelectContext(ctx, &pods, GetPodsByPlayerID, playerID); err != nil {
		return nil, fmt.Errorf("failed to get Pod records for player %d: %w", playerID, err)
	}
	if pods == nil {
		return []Pod{}, nil
	}
	return pods, nil
}

func (p *PodProvider) Add(ctx context.Context, name string) error {
	result, err := p.client.Db.ExecContext(ctx, InsertPod, name)
	if err != nil {
		return fmt.Errorf("failed to insert Pod record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Pod insert: got %d, expected 1", numAffected)
	}

	return nil
}

func (p *PodProvider) AddPlayerToPod(ctx context.Context, podID, playerID int) error {
	result, err := p.client.Db.ExecContext(ctx, InsertPlayerPod, podID, playerID)
	if err != nil {
		return fmt.Errorf("failed to insert PlayerPod record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by PlayerPod insert: got %d, expected 1", numAffected)
	}

	return nil
}
