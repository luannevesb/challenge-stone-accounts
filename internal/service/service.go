package service

import (
	"github.com/luannevesb/challenge-stone-accounts/internal/types"
	"net/http"
)

type Service struct {
	storageAccount types.StorageAccount
}

func NewService(storageAccount types.StorageAccount) *Service {
	return &Service{
		storageAccount: storageAccount,
	}
}
func GetAccount(w http.ResponseWriter, r *http.Request) {
	//TODO GetAccount
}

func CreateAccounts(w http.ResponseWriter, r *http.Request) {
	// TODO CreateAccounts
}

func GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	//TODO GetAllAccounts
}