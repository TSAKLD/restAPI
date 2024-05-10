package main

import (
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

	projRepo := repository.NewProjectRepository(db)
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	us := service.New(projRepo, userRepo, authRepo, taskRepo)

	hdr := api.NewHandler(us)
	mw := api.NewMiddleware(us)

	server := api.NewServer(hdr, cfg.HTTPPort, mw)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
