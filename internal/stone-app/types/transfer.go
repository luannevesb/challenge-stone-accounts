package types

type Transfer struct {
	Id                   string      `json:"id,omitempty"`
	AccountOriginId      string      `json:"account_origin_id"`
	AccountdestinationId string      `json:"account_destination_id"`
	Amount               float64     `json:"amount"`
	CreatedAt            interface{} `json:"created_at"`
}

type StorageTransfer interface {
	CreateTransfer(transfer *Transfer) error
	GetAllTranfers() ([]string, error)
}
