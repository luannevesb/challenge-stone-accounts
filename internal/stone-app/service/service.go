package service

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/types"
	"github.com/luannevesb/challenge-stone-accounts/pkg/helper"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"reflect"
	"time"
)

const (
	AttrinuteName                 = "name"
	AttrinuteCPF                  = "cpf"
	AttributeBallance             = "ballance"
	dateLayout                    = "2006-01-02"
	AttributeAccountOriginID      = "account_origin_id"
	AttributeAccountDestinationID = "account_destination_id"
	AttributeAmount               = "amount"
	ErroInesperado                = "Erro Inesperado"
	ErroNotFound                  = "Conta não encontrada"
	ErroOriginNotFound            = "Conta de origem não encontrada"
	ErroDestinationNotFound       = "Conta de destino não encontrada"
	ErroInsuficientBallance       = "Conta não possui saldo suficiente"
	ErrorResourceExist            = "Já existe conta criada com esse CPF"
)

var ValidateMessagesCreateAccount = govalidator.MapData{
	AttrinuteName:     {"required:O campo é obrigatório", "alpha: O campo tem formato inválido"},
	AttrinuteCPF:      {"required:O campo é obrigatório", "cpf: O CPF é inválido"},
	AttributeBallance: {"required:O campo é obrigatório", "float: O campo tem formato inválido"},
}

var ValidateRulesCreateAccount = map[string][]string{
	AttrinuteName:     {"required", "alpha"},
	AttrinuteCPF:      {"required", "cpf"},
	AttributeBallance: {"required", "float"},
}

var ValidateMessagesCreateTransfer = govalidator.MapData{
	AttributeAccountOriginID:      {"required:O campo é obrigatório"},
	AttributeAccountDestinationID: {"required:O campo é obrigatório"},
	AttributeAmount:               {"required:O campo é obrigatório","float: O campo tem formato inválido"},
}

var ValidateRulesCreateTransfer = map[string][]string{
	AttributeAccountOriginID:      {"required",},
	AttributeAccountDestinationID: {"required",},
	AttributeAmount:               {"required","float" },
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
	account := helper.ValidateJsonRequestAccount(w, r, ValidateRulesCreateAccount, ValidateMessagesCreateAccount)

	if account == nil {
		return
	}

	//Usando um base64 do CPF como ID já que não podem existir duas contas com o mesmo CPF
	id := base64.StdEncoding.EncodeToString([]byte(account.Cpf))

	var accountExistintent = types.Account{}

	_ = s.storageAccount.GetAccount(id, &accountExistintent)

	//Se já existe uma conta com esse CPF retorna 422
	if !reflect.ValueOf(accountExistintent).IsZero() {
		TrowError(w, http.StatusUnprocessableEntity, ErrorResourceExist)
		return
	}

	//Se não existe ele cria uma nova com o Created_at now()
	account.Id = id
	account.CreatedAt = time.Now().Format(dateLayout)

	err := s.storageAccount.CreateAccount(account)

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
	//Validação da request para os campos obrigatórios
	transfer := helper.ValidateJsonRequestTransfer(w, r, ValidateRulesCreateTransfer, ValidateMessagesCreateTransfer)

	if transfer == nil {
		return
	}

	//Verifica se a conta de origem existe
	accountOrigin := &types.Account{}
	err := s.storageAccount.GetAccount(transfer.AccountOriginId, accountOrigin)

	if err != nil {
		TrowError(w, http.StatusNotFound, ErroOriginNotFound)
		return
	}

	//Verifica se a conta de destino existe
	accountDestination := &types.Account{}
	err = s.storageAccount.GetAccount(transfer.AccountdestinationId, accountDestination)

	if err != nil {
		TrowError(w, http.StatusNotFound, ErroDestinationNotFound)
		return
	}

	//Verifica se o saldo da conta é menor que a transferência
	if accountOrigin.Ballance < transfer.Amount {
		TrowError(w, http.StatusUnprocessableEntity, ErroInsuficientBallance)
		return
	}

	//Atualiza o saldo na conta de origem
	accountOrigin.Ballance = accountOrigin.Ballance - transfer.Amount
	err = s.storageAccount.UpdateAccount(accountOrigin)

	if err != nil {
		TrowError(w, http.StatusInternalServerError, ErroInesperado)
		return
	}

	//Atualiza o saldo na conta de destino
	accountDestination.Ballance = accountDestination.Ballance + transfer.Amount
	err = s.storageAccount.UpdateAccount(accountDestination)

	if err != nil {
		TrowError(w, http.StatusInternalServerError, ErroInesperado)
		return
	}

	transfer.Id = uuid.New().String()
	transfer.CreatedAt = time.Now().Format(dateLayout)

	err = s.storageTransfer.CreateTransfer(transfer)

	if err != nil {
		TrowError(w, http.StatusInternalServerError, ErroInesperado)
		return
	}

	TrowSucess(w, types.SucessResponse{Sucess: true, Data: transfer})
}

func (s *Service) GetAllTransfers(w http.ResponseWriter, r *http.Request) {
	transfers, err := s.storageTransfer.GetAllTransfers()

	if err != nil {
		TrowError(w, http.StatusInternalServerError, ErroInesperado)
		return
	}

	var transfersJson []types.Transfer

	for _, transfer := range transfers {
		transferJson := &types.Transfer{}

		err = json.Unmarshal([]byte(transfer), transferJson)

		if err != nil {
			TrowError(w, http.StatusInternalServerError, ErroInesperado)
			return
		}

		transfersJson = append(transfersJson, *transferJson)
	}

	TrowSucess(w, types.SucessResponse{Sucess: true, Data: transfersJson})
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
