package types

type Account struct {
	Id        string      `json:"id,omitempty"`
	Name      string      `json:"name"`
	Cpf       string      `json:"cpf"`
	Ballance  float64     `json:"ballance"`
	CreatedAt interface{} `json:"created_at"`
}

type StorageAccount interface {
	CreateAccount(account *Account) error
	GetAccount(id string, account *Account) error
	GetAllAccounts() ([]string, error)
}
