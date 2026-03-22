package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const lookupFKName = `
	SELECT CONSTRAINT_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
	WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?
  	AND REFERENCED_TABLE_NAME IS NOT NULL LIMIT 1;
`

type Migration19 struct{}

func (m *Migration19) Upgrade(ctx context.Context, db *gorm.DB) error {
	steps := []struct{ table, column, name, ref, refCol string }{
		{"deck_commander", "deck_id", "fk_dc_deck", "deck", "id"},
		{"player_pod", "pod_id", "fk_player_pod_pod", "pod", "id"},
		{"player_pod", "player_id", "fk_player_pod_player", "player", "id"},
		{"user", "player_id", "fk_user_player", "player", "id"},
	}
	for _, s := range steps {
		drop := fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", s.table, s.name)
		add := fmt.Sprintf(
			"ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`) ON DELETE CASCADE",
			s.table, s.name, s.column, s.ref, s.refCol,
		)
		for _, stmt := range []string{drop, add} {
			if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
				return fmt.Errorf("query %q: %w", stmt, err)
			}
		}
	}

	// player_pod_role and pod_invite were created without explicit constraint names — look them up.
	dynamic := []struct{ table, column, ref, refCol string }{
		{"player_pod_role", "pod_id", "pod", "id"},
		{"player_pod_role", "player_id", "player", "id"},
		{"pod_invite", "pod_id", "pod", "id"},
	}
	for _, d := range dynamic {
		var constraintName string
		if err := db.WithContext(ctx).Raw(lookupFKName, d.table, d.column).Scan(&constraintName).Error; err != nil {
			return fmt.Errorf("lookup FK for %s.%s: %w", d.table, d.column, err)
		}
		if constraintName == "" {
			return fmt.Errorf("no FK constraint found for %s.%s", d.table, d.column)
		}
		drop := fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", d.table, constraintName)
		add := fmt.Sprintf(
			"ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`) ON DELETE CASCADE",
			d.table, constraintName, d.column, d.ref, d.refCol,
		)
		for _, stmt := range []string{drop, add} {
			if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
				return fmt.Errorf("query %q: %w", stmt, err)
			}
		}
	}
	return nil
}

func (m *Migration19) Downgrade(ctx context.Context, db *gorm.DB) error {
	steps := []struct{ table, column, name, ref, refCol string }{
		{"deck_commander", "deck_id", "fk_dc_deck", "deck", "id"},
		{"player_pod", "pod_id", "fk_player_pod_pod", "pod", "id"},
		{"player_pod", "player_id", "fk_player_pod_player", "player", "id"},
		{"user", "player_id", "fk_user_player", "player", "id"},
	}
	for _, s := range steps {
		drop := fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", s.table, s.name)
		add := fmt.Sprintf(
			"ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)",
			s.table, s.name, s.column, s.ref, s.refCol,
		)
		for _, stmt := range []string{drop, add} {
			if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
				return fmt.Errorf("query %q: %w", stmt, err)
			}
		}
	}

	dynamic := []struct{ table, column, ref, refCol string }{
		{"player_pod_role", "pod_id", "pod", "id"},
		{"player_pod_role", "player_id", "player", "id"},
		{"pod_invite", "pod_id", "pod", "id"},
	}
	for _, d := range dynamic {
		var constraintName string
		if err := db.WithContext(ctx).Raw(lookupFKName, d.table, d.column).Scan(&constraintName).Error; err != nil {
			return fmt.Errorf("lookup FK for %s.%s: %w", d.table, d.column, err)
		}
		if constraintName == "" {
			return fmt.Errorf("no FK constraint found for %s.%s", d.table, d.column)
		}
		drop := fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", d.table, constraintName)
		add := fmt.Sprintf(
			"ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)",
			d.table, constraintName, d.column, d.ref, d.refCol,
		)
		for _, stmt := range []string{drop, add} {
			if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
				return fmt.Errorf("query %q: %w", stmt, err)
			}
		}
	}
	return nil
}
