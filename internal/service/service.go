package service

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/luannevesb/challenge-stone-accounts/internal/types"
	"github.com/luannevesb/challenge-stone-accounts/pkg/helper"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"reflect"
	"time"
)

const (
	AttributeID        = "id"
	AttrinuteName      = "name"
	AttrinuteCPF       = "cpf"
	AttributeBallance  = "ballance"
	dateLayout         = "2006-01-02"
	AttributeCreatedAt = "created_at"
	ErroInesperado     = "Erro Inesperado"
	ErroNotFound       = "Conta não encontrada"
	ErrorResourceExist = "Já existe conta criada com esse CPF"
)

var ValidateMessagesCreateAccount = govalidator.MapData{
	AttributeID:        {"required: O campo é obrigatório", "numeric: O campo tem formato inválido"},
	AttrinuteName:      {"required: O campo é obrigatório", "alpha: O campo tem formato inválido"},
	AttrinuteCPF:       {"required: O campo é obrigatório", "cpf: O CPF é inválido"},
	AttributeBallance:  {"required: O campo é obrigatório", "float: O campo tem formato inválido"},
	AttributeCreatedAt: {"required: O campo é obrigatório", "date: A data é inválida"},
}

var ValidateRulesCreateAccount = map[string][]string{
	AttrinuteName:     {"required", "alpha"},
	AttrinuteCPF:      {"required", "cpf"},
	AttributeBallance: {"required", "float"},
}

func init() {
	helper.InitCustomRule()
}

type Service struct {
	storageAccount types.StorageAccount
}

func NewService(storageAccount types.StorageAccount) *Service {
	return &Service{
		storageAccount: storageAccount,
	}
}

func (s *Service) GetAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	account := &types.Account{}

	err := s.storageAccount.GetAccount(id, account)

	if err != nil {
		TrowError(w, http.StatusNotFound, ErroNotFound)
		return
	}

	TrowSucess(w, types.SucessResponse{Sucess: true, Data: account})
}

func (s *Service) CreateAccounts(w http.ResponseWriter, r *http.Request) {
	account := helper.ValidateJsonRequest(w, r, ValidateRulesCreateAccount, ValidateMessagesCreateAccount)

	if account == nil {
		return
	}

	id := base64.StdEncoding.EncodeToString([]byte(account.Cpf))

	var accountExistintent = types.Account{}

	s.storageAccount.GetAccount(id, &accountExistintent)

	if !reflect.ValueOf(accountExistintent).IsZero() {
		TrowError(w, http.StatusUnprocessableEntity, ErrorResourceExist)
		return
	}

	account.Id = id
	account.CreatedAt = time.Now().Format(dateLayout)

	err := s.storageAccount.CreateAccount(account)

	if err != nil {
		TrowError(w, http.StatusInternalServerError, ErroInesperado)
		return
	}

	TrowSucess(w, types.SucessResponse{Sucess: true, Data: account})
}

func (s *Service) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.storageAccount.GetAllAccounts()

	if err != nil {
		TrowError(w, http.StatusInternalServerError, ErroInesperado)
		return
	}
	var accountsJson []types.Account

	for _, account := range accounts {
		accountJson := &types.Account{}

		err = json.Unmarshal([]byte(account), accountJson)
		if err != nil {
			TrowError(w, http.StatusInternalServerError, ErroInesperado)
			return
		}

		accountsJson = append(accountsJson, *accountJson)
	}

	TrowSucess(w, types.SucessResponse{Sucess: true, Data: accountsJson})
}

func TrowError(w http.ResponseWriter, statusCode int, Error interface{}) {
	SetStatusCode(w, statusCode)
	SetJsonEncoder(w, types.ErrorResponse{Error: Error})
}

func TrowSucess(w http.ResponseWriter, Data interface{}) {
	SetStatusCode(w, http.StatusOK)
	SetJsonEncoder(w, Data)
}

func SetStatusCode(w http.ResponseWriter, status int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
}

func SetJsonEncoder(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
