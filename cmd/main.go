package main

import (
	"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/provider/http"
	"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/service"
	"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/storage"
	"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/storage/model"
)

func main() {
	//Cria instância do DB
	storage := storage.NewStorage()

	//Cria a instância do Model de Account
	storageAccount := &model.AccountStorage{DB: storage}

	//Cria a instância do Model de Transfer
	storageTransfer := &model.TransfersStorage{DB: storage}

	//Cria a instância de um novo Service passando os models
	service := service.NewService(storageAccount, storageTransfer)

	//Cria a instância das rotas do serviço
	http.InitRouter(service)
}
