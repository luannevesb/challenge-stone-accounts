package types

type MockStorageAccount struct {
	OnCreateAccount  func(account *Account) error
	OnGetAccount     func(id string, account *Account) error
	OnGetAllAccounts func() ([]string, error)
}

func (m *MockStorageAccount) CreateAccount(account *Account) error {
	return m.OnCreateAccount(account)
}

func (m *MockStorageAccount) GetAccount(id string, account *Account) error {
	return m.OnGetAccount(id, account)
}

func (m *MockStorageAccount) GetAllAccounts() ([]string, error) {
	return m.OnGetAllAccounts()
}
