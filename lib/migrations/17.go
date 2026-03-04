package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
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

func (m *Migration17) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createPlayerPodRoleTable); err != nil {
		return lib.NewDBError(createPlayerPodRoleTable, err)
	}
	return nil
}

func (m *Migration17) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropPlayerPodRoleTable); err != nil {
		return lib.NewDBError(dropPlayerPodRoleTable, err)
	}
	return nil
}
