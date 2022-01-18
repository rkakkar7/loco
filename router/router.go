package router

import (
	"loco/transaction"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/transactionservice/transaction/{transaction_id}", transaction.RegisterTransaction).Methods("PUT")
	router.HandleFunc("/transactionservice/transaction/{transaction_id}", transaction.GetTransactionByID).Methods("GET")
	router.HandleFunc("/transactionservice/types/{type}", transaction.GetTransactionsByType).Methods("GET")
	router.HandleFunc("/transactionservice/sum/{transaction_id}", transaction.GetTransactionSum).Methods("GET")
	router.StrictSlash(true)
	return router
}
