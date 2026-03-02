package main

import (
	"context"
	"log"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/migrations"
	"github.com/m-sharp/edh-tracker/lib/models"
	"github.com/m-sharp/edh-tracker/lib/repositories"
	"github.com/m-sharp/edh-tracker/lib/seeder"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := lib.NewConfig(lib.DBHost, lib.DBUsername, lib.DBPass, lib.DBPort)
	if err != nil {
		log.Fatalf("Error creating Config: %s", err.Error())
	}

	logger := lib.GetLogger(cfg)

	client, err := lib.NewDBClient(cfg, logger)
	if err != nil {
		log.Fatal("Error creating DB client", zap.Error(err))
	}

	if err = migrations.RunAll(ctx, client, logger); err != nil {
		log.Fatal("Failed to run DB migrations", zap.Error(err))
	}

	repos := models.NewRepositories(logger, client)

	if seed, _ := cfg.Get(lib.Seed); seed != "" {
		s := seeder.NewSeeder(logger, repos)
		if err = s.Run(ctx); err != nil {
			logger.Fatal("Seeder failed", zap.Error(err))
		}
	}

	repoLayer := repositories.New(logger, client)
	biz := business.NewBusiness(logger, repoLayer)

	apiRouter := NewApiRouter(cfg, logger, biz)
	server := lib.NewWebServer(cfg, logger, func(router *mux.Router) {
		// ToDo: Will need some auth for the app's connection eventually
		apiRouter.SetupRoutes(router)
	})

	if err := server.Serve(); err != nil {
		logger.Fatal("Server stopped listening", zap.Error(err))
	}
}
