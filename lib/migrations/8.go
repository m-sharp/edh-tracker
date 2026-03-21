package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createPod = `CREATE TABLE pod (
		id         INT          NOT NULL AUTO_INCREMENT,
		name       VARCHAR(255) NOT NULL,
		created_at DATETIME     NOT NULL DEFAULT NOW(),
		updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME     NULL,
		PRIMARY KEY (id)
	);`
	createPlayerPod = `CREATE TABLE player_pod (
		id         INT      NOT NULL AUTO_INCREMENT,
		pod_id     INT      NOT NULL,
		player_id  INT      NOT NULL,
		created_at DATETIME NOT NULL DEFAULT NOW(),
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME NULL,
		PRIMARY KEY (id),
		UNIQUE KEY uq_player_pod (pod_id, player_id),
		CONSTRAINT fk_player_pod_pod    FOREIGN KEY (pod_id)    REFERENCES pod (id),
		CONSTRAINT fk_player_pod_player FOREIGN KEY (player_id) REFERENCES player (id)
	);`

	dropPlayerPod = `DROP TABLE player_pod;`
	dropPod       = `DROP TABLE pod;`
)

type Migration8 struct{}

func (m *Migration8) Upgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{createPod, createPlayerPod} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}

func (m *Migration8) Downgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{dropPlayerPod, dropPod} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}
