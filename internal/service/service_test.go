package service_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/luannevesb/challenge-stone-accounts/internal/service"
	"github.com/luannevesb/challenge-stone-accounts/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	storageAccount *types.MockStorageAccount
	srv            *service.Service
)

//Teste Unitário da função CreateAccount do Service
func TestService_CreateAccount(t *testing.T) {
	srv := setup()

	t.Run("Teste CreateAccount Sucess", func(t *testing.T) {

		storageAccount.OnGetAccount = func(id string, account *types.Account) error {
			return nil
		}

		storageAccount.OnCreateAccount = func(account *types.Account) error {
			return nil
		}

		//Criando um FakeBody para request
		account := types.Account{Cpf: "70600657175", Name: "Luan", Ballance: 123.45}
		body, _ := json.Marshal(account)
		bytesBody := bytes.NewReader(body)

		//Criando nova fake request e recorder para teste
		mockedRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8080/accounts", bytesBody)
		recorder := httptest.NewRecorder()

		srv.CreateAccount(recorder, mockedRequest)

		//Verificação do comportamento de acordo com o cenário
		if recorder.Result().StatusCode != http.StatusOK {
			t.Errorf("Deveria retornar status 200; retornou %d", recorder.Result().StatusCode)
		}

	})

	t.Run("Teste CreateAccount CPF Exists", func(t *testing.T) {

		storageAccount.OnGetAccount = func(id string, account *types.Account) error {
			//Mockando um cenário onde já existe uma conta com esse CPF
			account.Ballance = 123.45
			account.Id = "xxx"
			account.Name = "Luan"
			account.Cpf = "04562342342"
			account.CreatedAt = time.Now()
			return nil
		}

		//Criando um FakeBody para request
		account := types.Account{Cpf: "70600657175", Name: "Luan", Ballance: 123.45}
		body, _ := json.Marshal(account)
		bytesBody := bytes.NewReader(body)

		//Criando nova fake request e recorder para teste
		mockedRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8080/accounts", bytesBody)
		recorder := httptest.NewRecorder()

		srv.CreateAccount(recorder, mockedRequest)

		if recorder.Result().StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Deveria retornar status 422; retornou %d", recorder.Result().StatusCode)
		}

	})

	t.Run("Teste CreateAccount ErroInesperado", func(t *testing.T) {

		storageAccount.OnGetAccount = func(id string, account *types.Account) error {
			return nil
		}

		storageAccount.OnCreateAccount = func(account *types.Account) error {
			return errors.New("Erro Inesperado")
		}

		//Criando um FakeBody para request
		account := types.Account{Cpf: "70600657175", Name: "Luan", Ballance: 123.45}
		body, _ := json.Marshal(account)
		bytesBody := bytes.NewReader(body)

		//Criando nova fake request e recorder para teste
		mockedRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8080/accounts", bytesBody)
		recorder := httptest.NewRecorder()

		srv.CreateAccount(recorder, mockedRequest)

		if recorder.Result().StatusCode != http.StatusInternalServerError {
			t.Errorf("Deveria retornar status 500; retornou %d", recorder.Result().StatusCode)
		}

	})

}

//Teste Unitário da função GetAccount do Service
func TestService_GetAccount(t *testing.T) {
	srv = setup()
	t.Run("Teste GetAccount Sucess", func(t *testing.T) {

		storageAccount.OnGetAccount = func(id string, account *types.Account) error {
			return nil
		}

		//Criando nova fake request e recorder para teste
		mockedRequest := httptest.NewRequest(http.MethodGet, "http://localhost:8080/accounts/", nil)
		recorder := httptest.NewRecorder()

		srv.GetAccount(recorder, mockedRequest)

		if recorder.Result().StatusCode != http.StatusOK {
			t.Errorf("Deveria retornar status 200; retornou %d", recorder.Result().StatusCode)
		}

	})

	t.Run("Teste GetAccount ErroNotFound", func(t *testing.T) {

		storageAccount.OnGetAccount = func(id string, account *types.Account) error {
			return errors.New("Objeto não encontrado")
		}

		//Criando nova fake request e recorder para teste
		mockedRequest := httptest.NewRequest(http.MethodGet, "http://localhost:8080/accounts/", nil)
		recorder := httptest.NewRecorder()

		srv.GetAccount(recorder, mockedRequest)

		if recorder.Result().StatusCode != http.StatusNotFound {
			t.Errorf("Deveria retornar status 404; retornou %d", recorder.Result().StatusCode)
		}

	})
}

//Teste Unitário da função GetAllAccounts do Service
func TestService_GetAllAccounts(t *testing.T) {
	srv = setup()
	t.Run("Teste GetAllAccounts Sucess", func(t *testing.T) {

		storageAccount.OnGetAllAccounts = func() ([]string, error) {
			//Mockando um retorno de sucesso do Storage
			values := []string{
				`{"id": "MDUyMjAxNjAxNDE=","name": "Ketellen","cpf": "05220160141","ballance": 128.45,"created_at": "2020-02-15"}`,
			}
			return values, nil
		}

		//Criando nova fake request e recorder para teste
		mockedRequest := httptest.NewRequest(http.MethodGet, "http://localhost:8080/accounts/", nil)
		recorder := httptest.NewRecorder()

		srv.GetAllAccounts(recorder, mockedRequest)

		if recorder.Result().StatusCode != http.StatusOK {
			t.Errorf("Deveria retornar status 200; retornou %d", recorder.Result().StatusCode)
		}

	})

	t.Run("Teste GetAccount ErroInesperado", func(t *testing.T) {
		storageAccount.OnGetAllAccounts = func() (strings []string, err error) {
			//Mockando erro retornado do Storage
			return []string{}, errors.New("Erro Inesperado")
		}

		//Criando nova fake request e recorder para teste
		mockedRequest := httptest.NewRequest(http.MethodGet, "http://localhost:8080/accounts/", nil)
		recorder := httptest.NewRecorder()

		srv.GetAllAccounts(recorder, mockedRequest)

		if recorder.Result().StatusCode != http.StatusInternalServerError {
			t.Errorf("Deveria retornar status 500; retornou %d", recorder.Result().StatusCode)
		}

	})
}

//Cria novo MOCK do Model de Account e retorna instância de Service
func setup() *service.Service {
	storageAccount = &types.MockStorageAccount{}
	return service.NewService(storageAccount)

}
