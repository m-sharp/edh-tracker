package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
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

func (m *Migration18) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createPodInviteTable); err != nil {
		return lib.NewDBError(createPodInviteTable, err)
	}
	return nil
}

func (m *Migration18) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropPodInviteTable); err != nil {
		return lib.NewDBError(dropPodInviteTable, err)
	}
	return nil
}
