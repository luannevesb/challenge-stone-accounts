package storage

import (
	"github.com/nanobox-io/golang-scribble"
)

func NewStorage() *scribble.Driver {
	//Cria novo JSON Database (SCRIBBLE)
	db, _ := scribble.New("./accounts", nil)

	return db
}
