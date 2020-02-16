package model

import (
"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/types"
scribble "github.com/nanobox-io/golang-scribble"
)

type TransfersStorage struct {
	DB *scribble.Driver
}

//Função responsável por CRIAR Transfers no DB
func (s *TransfersStorage) CreateTransfer(transfer *types.Transfer) error {
	return s.DB.Write("transfers", transfer.Id, transfer)
}

//Função responsável por BUSCAR TODAS AS Transfers no DB
func (s *TransfersStorage) GetAllTranfers() ([]string, error) {
	return s.DB.ReadAll("transfers")
}
