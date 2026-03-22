package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	addPlayerIndexes     = `ALTER TABLE player ADD INDEX idx_player_name (name), ADD INDEX idx_player_deleted_at (deleted_at);`
	addDeckIndexes       = `ALTER TABLE deck ADD INDEX idx_deck_deleted_at (deleted_at), ADD INDEX idx_deck_retired (retired), ADD INDEX idx_deck_commander (commander), ADD INDEX idx_deck_player_commander (player_id, commander);`
	addGameIndexes       = `ALTER TABLE game ADD INDEX idx_game_deleted_at (deleted_at);`
	addGameResultIndexes = `ALTER TABLE game_result ADD INDEX idx_game_result_deleted_at (deleted_at);`
	addPodIndexes        = `ALTER TABLE pod ADD INDEX idx_pod_name (name), ADD INDEX idx_pod_deleted_at (deleted_at);`
	addPlayerPodIndexes  = `ALTER TABLE player_pod ADD INDEX idx_player_pod_deleted_at (deleted_at);`
	addUserIndexes       = `ALTER TABLE user ADD INDEX idx_user_deleted_at (deleted_at);`

	dropPlayerIndexes     = `ALTER TABLE player DROP INDEX idx_player_name, DROP INDEX idx_player_deleted_at;`
	dropDeckIndexes       = `ALTER TABLE deck DROP INDEX idx_deck_deleted_at, DROP INDEX idx_deck_retired, DROP INDEX idx_deck_commander, DROP INDEX idx_deck_player_commander;`
	dropGameIndexes       = `ALTER TABLE game DROP INDEX idx_game_deleted_at;`
	dropGameResultIndexes = `ALTER TABLE game_result DROP INDEX idx_game_result_deleted_at;`
	dropPodIndexes        = `ALTER TABLE pod DROP INDEX idx_pod_name, DROP INDEX idx_pod_deleted_at;`
	dropPlayerPodIndexes  = `ALTER TABLE player_pod DROP INDEX idx_player_pod_deleted_at;`
	dropUserIndexes       = `ALTER TABLE user DROP INDEX idx_user_deleted_at;`
)

type Migration10 struct{}

func (m *Migration10) Upgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{
		addPlayerIndexes,
		addDeckIndexes,
		addGameIndexes,
		addGameResultIndexes,
		addPodIndexes,
		addPlayerPodIndexes,
		addUserIndexes,
	} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}

func (m *Migration10) Downgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{
		dropPlayerIndexes,
		dropDeckIndexes,
		dropGameIndexes,
		dropGameResultIndexes,
		dropPodIndexes,
		dropPlayerPodIndexes,
		dropUserIndexes,
	} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}
