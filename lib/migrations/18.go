package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createPodInviteTable = `CREATE TABLE pod_invite (
    id                   INT AUTO_INCREMENT PRIMARY KEY,
    pod_id               INT NOT NULL,
    invite_code          VARCHAR(36) NOT NULL UNIQUE,
    created_by_player_id INT NOT NULL,
    expires_at           TIMESTAMP NULL,
    used_count           INT NOT NULL DEFAULT 0,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at           DATETIME NULL,
    INDEX idx_pi_pod_id      (pod_id),
    INDEX idx_pi_invite_code (invite_code),
    INDEX idx_pi_deleted_at  (deleted_at),
    FOREIGN KEY (pod_id)               REFERENCES pod(id),
    FOREIGN KEY (created_by_player_id) REFERENCES player(id)
);`

	dropPodInviteTable = `DROP TABLE IF EXISTS pod_invite;`
)

type Migration18 struct{}

func (m *Migration18) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(createPodInviteTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", createPodInviteTable, err)
	}
	return nil
}

func (m *Migration18) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(dropPodInviteTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", dropPodInviteTable, err)
	}
	return nil
}
