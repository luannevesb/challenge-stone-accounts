package http

import (
	"github.com/gorilla/mux"
	"github.com/luannevesb/challenge-stone-accounts/internal/stone-app/service"
	"log"
	"net/http"
)

func InitRouter(service *service.Service) {
	//Cria a instância de um novo ROUTER
	router := mux.NewRouter()

	//Inicia as rotas de accounts e informa qual método interno vai receber a REQUEST
	router.HandleFunc("/accounts", service.GetAllAccounts).Methods("GET")
	router.HandleFunc("/accounts/{id}", service.GetAccount).Methods("GET")
	router.HandleFunc("/accounts", service.CreateAccount).Methods("POST")

	//Inicia as rotas de transfers e informa qual método interno vai receber a REQUEST
	router.HandleFunc("/transfers", service.GetAllTransfers).Methods("GET")
	router.HandleFunc("/transfers", service.CreateTransfer).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
