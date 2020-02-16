package main

import (
	"github.com/luannevesb/challenge-stone-accounts/internal/provider/http"
	"github.com/luannevesb/challenge-stone-accounts/internal/service"
	"github.com/luannevesb/challenge-stone-accounts/internal/storage"
	"github.com/luannevesb/challenge-stone-accounts/internal/storage/model"
)

func main() {
	//Cria instância do DB
	storage := storage.NewStorage()

	//Cria a instância do Model de Account
	storageAccout := &model.AccountStorage{DB: storage}

	//Cria a instância de um novo Service passando os models
	service := service.NewService(storageAccout)

	//Cria a instância das rotas do serviço
	http.InitRouter(service)
}
