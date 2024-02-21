package controller

import "palco-planner-api/api/repository"

type Controller struct {
	storage repository.Storage
}

func NewController(storage repository.Storage) Controller {
	return Controller{storage}
}
