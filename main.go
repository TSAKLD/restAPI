package main

import (
	_ "github.com/lib/pq"
	"log"
	"restAPI/api"
	"restAPI/bootstrap"
	"restAPI/repository"
	"restAPI/service"
)

func main() {
	cfg, err := bootstrap.NewConfig()
	if err != nil {
		log.Fatal("Problem with config load: ", err)
	}

	errorList := cfg.Validate()
	if errorList != nil {
		log.Fatal("Problem with config validation: ", errorList)
	}

	db, err := bootstrap.DBConnect(cfg)
	if err != nil {
		log.Fatal("Problem with Database connection: ", err)
	}
	defer db.Close()

	log.Println("postgres DB connection status: OK")

	repo := repository.New(db)

	us := service.New(repo)

	hdr := api.NewHandler(us)

	server := api.NewServer(hdr, cfg.HTTPPort)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
