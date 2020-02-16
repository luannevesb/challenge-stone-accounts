package types

type MockStorageAccount struct {
	OnCreateAccount  func(account *Account) error
	OnGetAccount     func(id string, account *Account) error
	OnGetAllAccounts func() ([]string, error)
	OnUpdateAccount  func(account *Account) error
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

func (m *MockStorageAccount) UpdateAccount(account *Account) error {
	return m.OnUpdateAccount(account)
}
