package http

import (
	"github.com/gorilla/mux"
	"github.com/luannevesb/challenge-stone-accounts/internal/service"
	"log"
	"net/http"
)

func InitRouter(service *service.Service) {
	//Cria a instância de um novo ROUTER
	router := mux.NewRouter()

	//Inicia as rotas de accounts e informa qual método interno vai receber qual REQUEST
	router.HandleFunc("/accounts", service.GetAllAccounts).Methods("GET")
	router.HandleFunc("/accounts/{id}", service.GetAccount).Methods("GET")
	router.HandleFunc("/accounts", service.CreateAccount).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
