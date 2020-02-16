package model

import (
	"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/types"
	scribble "github.com/nanobox-io/golang-scribble"
)

type AccountStorage struct {
	DB *scribble.Driver
}

//Função responsável por CRIAR Account no DB
func (s *AccountStorage) CreateAccount(account *types.Account) error {
	return s.DB.Write("accounts", account.Id, account)
}

//Função responsável por BUSCAR Account no DB
func (s *AccountStorage) GetAccount(id string, account *types.Account) error {
	return s.DB.Read("accounts", id, account)
}

//Função responsável por BUSCAR TODAS AS Account no DB
func (s *AccountStorage) GetAllAccounts() ([]string, error) {
	return s.DB.ReadAll("accounts")
}

//Função responsável por ATUALIZAR uma Account no DB
func (s *AccountStorage) UpdateAccount(account *types.Account) error {
	return s.DB.Write("accounts", account.Id, account)
}


