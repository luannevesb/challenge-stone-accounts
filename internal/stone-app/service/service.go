package service

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/types"
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
	//Cria a regra customizada de CPF
	helper.InitCustomRule()
}

type Service struct {
	storageAccount  types.StorageAccount
	storageTransfer types.StorageTransfer
}

func NewService(storageAccount types.StorageAccount, storageTransfer types.StorageTransfer) *Service {
	return &Service{
		storageAccount:  storageAccount,
		storageTransfer: storageTransfer,
	}
}

//Retorna as informações da Account de acordo com o ID ou retorna 404
func (s *Service) GetAccount(w http.ResponseWriter, r *http.Request) {
	//Recebe os parâmetros da request e seleciona o ID
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

//Cria uma nova Account e retorna se existe account com esse CPF
func (s *Service) CreateAccount(w http.ResponseWriter, r *http.Request) {
	//Validação da request para os campos obrigatórios e o CPF
	account := helper.ValidateJsonRequest(w, r, ValidateRulesCreateAccount, ValidateMessagesCreateAccount)

	if account == nil {
		return
	}

	//Usando um base64 do CPF como ID já que não podem existir duas contas com o mesmo CPF
	id := base64.StdEncoding.EncodeToString([]byte(account.Cpf))

	var accountExistintent = types.Account{}

	err := s.storageAccount.GetAccount(id, &accountExistintent)

	//Se já existe uma conta com esse CPF retorna 422
	if !reflect.ValueOf(accountExistintent).IsZero() {
		TrowError(w, http.StatusUnprocessableEntity, ErrorResourceExist)
		return
	}

	if err != nil {
		TrowError(w, http.StatusInternalServerError, ErroInesperado)
		return
	}

	//Se não existe ele cria uma nova com o Created_at now()
	account.Id = id
	account.CreatedAt = time.Now().Format(dateLayout)

	err = s.storageAccount.CreateAccount(account)

	if err != nil {
		TrowError(w, http.StatusInternalServerError, ErroInesperado)
		return
	}

	TrowSucess(w, types.SucessResponse{Sucess: true, Data: account})
}

//Retorna as informações de todas as contas se não existir retorna []
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

func (s *Service) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	//TODO CreateTransfer
}

func (s *Service) GetAllTransfers(w http.ResponseWriter, r *http.Request) {
	//TODO GetAllTransfers
}

//Retorna um Erro de acordo com o formato passado e com o StatusCode informado
func TrowError(w http.ResponseWriter, statusCode int, Error interface{}) {
	SetStatusCode(w, statusCode)
	SetJsonEncoder(w, types.ErrorResponse{Error: Error})
}

//Retorna Sucesso com o Data informado
func TrowSucess(w http.ResponseWriter, Data interface{}) {
	SetStatusCode(w, http.StatusOK)
	SetJsonEncoder(w, Data)
}

//Seta o StatusCode no Writer
func SetStatusCode(w http.ResponseWriter, status int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
}

//Seta um Novo Encoder no Writer
func SetJsonEncoder(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
