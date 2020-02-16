package types

type MockStorageTransfer struct {
	OnCreateTransfer func(transfer *Transfer) error
	OnGetAllTransfers func() ([]string, error)
}

func (m *MockStorageTransfer) CreateTransfer(transfer *Transfer) error {
	return m.OnCreateTransfer(transfer)
}

func (m *MockStorageTransfer) GetAllTransfers() ([]string, error) {
	return m.OnGetAllTransfers()
}
