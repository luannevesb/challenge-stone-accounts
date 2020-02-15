package model

import (
	"github.com/luannevesb/challenge-stone-accounts/internal/types"
	scribble "github.com/nanobox-io/golang-scribble"
)

type AccountStorage struct {
	DB *scribble.Driver
}

func (s *AccountStorage)CreateAccount(account *types.Account) error {
	return s.DB.Write("accounts", account.Id, account)
}

func (s *AccountStorage)GetAccount(id string, account *types.Account) error {
	return s.DB.Read("accounts", id, account)
}

func (s *AccountStorage)GetAllAccounts() ([]string, error) {
	return s.DB.ReadAll("accounts")
}