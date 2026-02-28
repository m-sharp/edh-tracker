package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	createUserRole = `CREATE TABLE user_role (
		id         INT         NOT NULL AUTO_INCREMENT,
		created_at DATETIME    NOT NULL DEFAULT NOW(),
		updated_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME    NULL,
		name       VARCHAR(50) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY uq_user_role_name (name)
	);`
	seedUserRoles = `INSERT INTO user_role (name) VALUES ('admin'), ('player');`
	createUser    = `CREATE TABLE user (
		id             INT          NOT NULL AUTO_INCREMENT,
		player_id      INT          NOT NULL,
		role_id        INT          NOT NULL,
		oauth_provider VARCHAR(50)  NULL,
		oauth_subject  VARCHAR(255) NULL,
		email          VARCHAR(255) NULL,
		display_name   VARCHAR(255) NULL,
		avatar_url     TEXT         NULL,
		created_at     DATETIME     NOT NULL DEFAULT NOW(),
		updated_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at     DATETIME     NULL,
		PRIMARY KEY (id),
		UNIQUE KEY uq_user_player_id (player_id),
		UNIQUE KEY uq_user_oauth (oauth_provider, oauth_subject),
		CONSTRAINT fk_user_player FOREIGN KEY (player_id) REFERENCES player (id),
		CONSTRAINT fk_user_role   FOREIGN KEY (role_id)   REFERENCES user_role (id)
	);`

	dropUser     = `DROP TABLE user;`
	dropUserRole = `DROP TABLE user_role;`
)

type Migration7 struct{}

func (m *Migration7) Upgrade(ctx context.Context, client *lib.DBClient) error {
	for _, stmt := range []string{createUserRole, seedUserRoles, createUser} {
		if _, err := client.Db.ExecContext(ctx, stmt); err != nil {
			return lib.NewDBError(stmt, err)
		}
	}
	return nil
}

func (m *Migration7) Downgrade(ctx context.Context, client *lib.DBClient) error {
	for _, stmt := range []string{dropUser, dropUserRole} {
		if _, err := client.Db.ExecContext(ctx, stmt); err != nil {
			return lib.NewDBError(stmt, err)
		}
	}
	return nil
}
