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

	userServ := service.NewUserService(userRepo, authRepo)
	authServ := service.NewAuthService(authRepo)
	projServ := service.NewProjectRepository(projRepo, taskRepo)

	taskHandler := api.NewTaskHandler(projServ)
	projectHandler := api.NewProjectHandler(projServ)
	userHandler := api.NewUserHandler(userServ)
	authHandler := api.NewAuthHandler(authServ)

	mw := api.NewMiddleware(projServ, authServ, userServ)

	server := api.NewServer(taskHandler, projectHandler, userHandler, authHandler, cfg.HTTPPort, mw)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
