package main

import (
	"log"

	"install-scripts-controller/internal/config"
	"install-scripts-controller/internal/database"
	httpsrv "install-scripts-controller/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Open(cfg.Database.Path)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	srv := httpsrv.New(cfg, db, "web/dist")
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
