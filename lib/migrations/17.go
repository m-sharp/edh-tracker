package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createPlayerPodRoleTable = `CREATE TABLE player_pod_role (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    pod_id      INT NOT NULL,
    player_id   INT NOT NULL,
    role        ENUM('manager', 'member') NOT NULL DEFAULT 'member',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME NULL,
    UNIQUE KEY uq_ppr (pod_id, player_id),
    INDEX idx_ppr_pod_id    (pod_id),
    INDEX idx_ppr_player_id (player_id),
    INDEX idx_ppr_deleted_at (deleted_at),
    FOREIGN KEY (pod_id)    REFERENCES pod(id),
    FOREIGN KEY (player_id) REFERENCES player(id)
);`

	dropPlayerPodRoleTable = `DROP TABLE IF EXISTS player_pod_role;`
)

type Migration17 struct{}

func (m *Migration17) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(createPlayerPodRoleTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", createPlayerPodRoleTable, err)
	}
	return nil
}

func (m *Migration17) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(dropPlayerPodRoleTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", dropPlayerPodRoleTable, err)
	}
	return nil
}
