package pod

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getAllPods        = `SELECT id, name, created_at, updated_at, deleted_at FROM pod WHERE deleted_at IS NULL;`
	getPodByID        = `SELECT id, name, created_at, updated_at, deleted_at FROM pod WHERE id = ? AND deleted_at IS NULL;`
	getPodByName      = `SELECT id, name, created_at, updated_at, deleted_at FROM pod WHERE name = ? AND deleted_at IS NULL LIMIT 1;`
	getPodsByPlayerID = `SELECT pod.id, pod.name, pod.created_at, pod.updated_at, pod.deleted_at
						   FROM (pod INNER JOIN player_pod ON pod.id = player_pod.pod_id)
						  WHERE player_pod.player_id = ?
						    AND pod.deleted_at IS NULL
						    AND player_pod.deleted_at IS NULL;`
	getPodIDsByPlayerID  = `SELECT pod_id FROM player_pod WHERE player_id = ? AND deleted_at IS NULL;`
	getPlayerIDsByPodID  = `SELECT player_id FROM player_pod WHERE pod_id = ? AND deleted_at IS NULL;`
	insertPod            = `INSERT INTO pod (name) VALUES (?);`
	insertPlayerPod      = `INSERT INTO player_pod (pod_id, player_id) VALUES (?, ?);`
	softDeletePod        = `UPDATE pod SET deleted_at = NOW() WHERE id = ?;`
	updatePodName        = `UPDATE pod SET name = ? WHERE id = ? AND deleted_at IS NULL;`
	softDeletePlayerPod  = `UPDATE player_pod SET deleted_at = NOW() WHERE pod_id = ? AND player_id = ? AND deleted_at IS NULL;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
	var pods []Model
	if err := r.client.Db.SelectContext(ctx, &pods, getAllPods); err != nil {
		return nil, fmt.Errorf("failed to get Pod records: %w", err)
	}
	if pods == nil {
		return []Model{}, nil
	}
	return pods, nil
}

func (r *Repository) GetByID(ctx context.Context, podID int) (*Model, error) {
	var pods []Model
	if err := r.client.Db.SelectContext(ctx, &pods, getPodByID, podID); err != nil {
		return nil, fmt.Errorf("failed to get Pod record for id %d: %w", podID, err)
	}
	if len(pods) == 0 {
		return nil, nil
	}
	return &pods[0], nil
}

func (r *Repository) GetByPlayerID(ctx context.Context, playerID int) ([]Model, error) {
	var pods []Model
	if err := r.client.Db.SelectContext(ctx, &pods, getPodsByPlayerID, playerID); err != nil {
		return nil, fmt.Errorf("failed to get Pod records for player %d: %w", playerID, err)
	}
	if pods == nil {
		return []Model{}, nil
	}
	return pods, nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*Model, error) {
	var pods []Model
	if err := r.client.Db.SelectContext(ctx, &pods, getPodByName, name); err != nil {
		return nil, fmt.Errorf("failed to get Pod record for name %q: %w", name, err)
	}
	if len(pods) == 0 {
		return nil, nil
	}
	return &pods[0], nil
}

func (r *Repository) GetIDsByPlayerID(ctx context.Context, playerID int) ([]int, error) {
	var ids []int
	if err := r.client.Db.SelectContext(ctx, &ids, getPodIDsByPlayerID, playerID); err != nil {
		return nil, fmt.Errorf("failed to get Pod IDs for player %d: %w", playerID, err)
	}
	if ids == nil {
		return []int{}, nil
	}
	return ids, nil
}

func (r *Repository) Add(ctx context.Context, name string) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertPod, name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert Pod record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return 0, fmt.Errorf("unexpected number of rows affected by Pod insert: got %d, expected 1", numAffected)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new Pod: %w", err)
	}

	return int(id), nil
}

func (r *Repository) BulkAddPlayers(ctx context.Context, podID int, playerIDs []int) error {
	if len(playerIDs) == 0 {
		return nil
	}

	query := "INSERT INTO player_pod (pod_id, player_id) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?),", len(playerIDs)), ",")
	args := make([]interface{}, 0, len(playerIDs)*2)
	for _, id := range playerIDs {
		args = append(args, podID, id)
	}
	if _, err := r.client.Db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk insert PlayerPod records: %w", err)
	}
	return nil
}

func (r *Repository) AddPlayerToPod(ctx context.Context, podID, playerID int) error {
	result, err := r.client.Db.ExecContext(ctx, insertPlayerPod, podID, playerID)
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

func (r *Repository) SoftDelete(ctx context.Context, podID int) error {
	if _, err := r.client.Db.ExecContext(ctx, softDeletePod, podID); err != nil {
		return fmt.Errorf("failed to soft delete pod %d: %w", podID, err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, podID int, name string) error {
	if _, err := r.client.Db.ExecContext(ctx, updatePodName, name, podID); err != nil {
		return fmt.Errorf("failed to update pod %d name: %w", podID, err)
	}
	return nil
}

func (r *Repository) RemovePlayer(ctx context.Context, podID, playerID int) error {
	if _, err := r.client.Db.ExecContext(ctx, softDeletePlayerPod, podID, playerID); err != nil {
		return fmt.Errorf("failed to remove player %d from pod %d: %w", playerID, podID, err)
	}
	return nil
}

func (r *Repository) GetPlayerIDs(ctx context.Context, podID int) ([]int, error) {
	var ids []int
	if err := r.client.Db.SelectContext(ctx, &ids, getPlayerIDsByPodID, podID); err != nil {
		return nil, fmt.Errorf("failed to get player IDs for pod %d: %w", podID, err)
	}
	if ids == nil {
		return []int{}, nil
	}
	return ids, nil
}
