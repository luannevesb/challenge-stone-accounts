package main

import (
	"github.com/luannevesb/challenge-stone-accounts/internal/provider/http"
	"github.com/luannevesb/challenge-stone-accounts/internal/service"
	"github.com/luannevesb/challenge-stone-accounts/internal/storage"
	"github.com/luannevesb/challenge-stone-accounts/internal/storage/model"
)

func main () {
	storage := storage.NewStorage()
	storageAccout := &model.AccountStorage{DB: storage}
	service :=  service.NewService(storageAccout)
	http.InitRouter(service)
}
