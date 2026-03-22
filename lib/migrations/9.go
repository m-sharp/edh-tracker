package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	addPodIDToGame    = `ALTER TABLE game ADD COLUMN pod_id INT NOT NULL, ADD CONSTRAINT fk_game_pod FOREIGN KEY (pod_id) REFERENCES pod (id);`
	dropPodIDFromGame = `ALTER TABLE game DROP FOREIGN KEY fk_game_pod, DROP COLUMN pod_id;`
)

type Migration9 struct{}

func (m *Migration9) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(addPodIDToGame).Error; err != nil {
		return fmt.Errorf("query %q: %w", addPodIDToGame, err)
	}
	return nil
}

func (m *Migration9) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(dropPodIDFromGame).Error; err != nil {
		return fmt.Errorf("query %q: %w", dropPodIDFromGame, err)
	}
	return nil
}
