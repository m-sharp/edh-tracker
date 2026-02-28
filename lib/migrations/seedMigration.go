package migrations

import (
	"context"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/migrations/seeder"
)

// ToDo: This should probably live on its own path entirely, but I'll just run it as the last migration for now

type SeedMigration struct {
	seeder *seeder.Seeder
}

func NewSeedMigration(log *zap.Logger, client *lib.DBClient) *SeedMigration {
	return &SeedMigration{
		seeder: seeder.NewSeeder(log, client),
	}
}

func (s SeedMigration) Upgrade(ctx context.Context, _ *lib.DBClient) error {
	return s.seeder.Run(ctx)
}

func (s SeedMigration) Downgrade(ctx context.Context, client *lib.DBClient) error {
	// There's no going back now...
	return nil
}

func (s SeedMigration) RecordMigration() bool { return false }
