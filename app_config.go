package main

import (
	"palco-planner-api/api"
	"palco-planner-api/controller"
	repository_firebase_storage "palco-planner-api/repository/firebase_storage"
)

type AppConfig struct {
	server api.Server
}

func NewAppConfig() AppConfig {

	repositoryFirebaseStorage := repository_firebase_storage.NewFirebaseStorage()
	controller := controller.NewController(repositoryFirebaseStorage)
	server := api.NewServer(controller)
	return AppConfig{
		server: server,
	}
}
