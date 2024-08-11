package pkg

import (
	"transaction-service/internal"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/transactionservice/transaction/{id:[0-9]+}", internal.CreateTransactionHandler).Methods("PUT")
	r.HandleFunc("/transactionservice/transaction/{id:[0-9]+}", internal.GetTransactionHandler).Methods("GET")
	r.HandleFunc("/transactionservice/types/{type}", internal.GetTransactionsByTypeHandler).Methods("GET")
	r.HandleFunc("/transactionservice/sum/{id:[0-9]+}", internal.GetSumHandler).Methods("GET")
	return r
}
