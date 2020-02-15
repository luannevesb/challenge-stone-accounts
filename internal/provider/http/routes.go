package http

import (
	"github.com/gorilla/mux"
	"github.com/luannevesb/challenge-stone-accounts/internal/service"
	"log"
	"net/http"
)

func InitRouter (service *service.Service) {
	//Inicia o a dependência MUX
	router := mux.NewRouter()

	//Inicia as rotas de accounts e informa qual método vai receber qual chamada
	router.HandleFunc("/accounts", service.GetAllAccounts).Methods("GET")
	router.HandleFunc("/accounts/{id}", service.GetAccount).Methods("GET")
	router.HandleFunc("/accounts", service.CreateAccounts).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}