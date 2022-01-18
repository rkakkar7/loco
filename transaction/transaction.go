package transaction

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"loco/api"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "raghav"
	password = "1234"
	dbname   = "db_transaction"
)

var db *sql.DB

var sqlInsert string = `INSERT INTO transaction (transaction_id, amount, type, parent_id) VALUES ($1, $2, $3, $4)`
var sqlGetByTxnID string = `SELECT amount, type, parent_id from transaction WHERE transaction_id = $1`
var sqlGetByType string = `SELECT transaction_id from transaction WHERE type = $1`
var sqlGetSum string = `WITH RECURSIVE q AS (
    SELECT transaction_id, parent_id, amount FROM transaction WHERE transaction_id = $1
  	UNION
    SELECT p.transaction_id, p.parent_id, p.amount
    FROM transaction p JOIN q ON p.parent_id = q.transaction_id
  )
SELECT sum(amount) from q;`

func init() {
	initDB()
}

func initDB() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func RegisterTransaction(w http.ResponseWriter, r *http.Request) {
	log.Info("RegisterTransaction Request received")
	params := mux.Vars(r)
	var txnID string

	var registerTransaction api.Transaction
	err := json.NewDecoder(r.Body).Decode(&registerTransaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("RegisterTransaction: json.NewDecoder %v", err)
		return
	}

	var ok bool
	if txnID, ok = params["transaction_id"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("transaction_id not present"))
		log.Errorf("RegisterTransaction: transaction_id not present")
		return
	}

	if registerTransaction.ParentID == nil {
		*registerTransaction.ParentID, _ = strconv.ParseInt(txnID, 10, 64)
	}

	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("RegisterTransaction: db.BeginTx %w", err)
		return
	}

	defer tx.Rollback()
	_, err = tx.ExecContext(r.Context(), sqlInsert, txnID, registerTransaction.Amount, registerTransaction.Type, registerTransaction.ParentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("RegisterTransaction: tx.ExecContext %w", err)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("RegisterTransaction: tx.Commit %w", err)
		return
	}

	log.Println("RegisterTransaction Successfull")
	w.WriteHeader(http.StatusOK)
}

func GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var txnID string
	var ok bool
	if txnID, ok = params["transaction_id"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("transaction_id not present"))
		log.Error("GetTransactionByID: transaction_id not present")
		return
	}

	row := db.QueryRowContext(r.Context(), sqlGetByTxnID, txnID)
	if row.Err() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("GetTransactionByID: db.QueryRowContext: %v", row.Err())
		return
	}

	var t api.Transaction
	err := row.Scan(&t.Amount, &t.Type, &t.ParentID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		log.Infof("GetTransactionByID: row.Scan sql.ErrNoRows: %v", row.Err())
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("GetTransactionByID: row.Scan: %v", err)
		return
	}

	json.NewEncoder(w).Encode(t)
}

func GetTransactionsByType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tag string
	var ok bool
	if tag, ok = params["type"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("GetTransactionsByType: type not present"))
		return
	}

	rows, err := db.QueryContext(r.Context(), sqlGetByType, tag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("GetTransactionsByType: db.ExecContext: %v", err)
		return
	}

	txnIDs := make([]int64, 0)
	for rows.Next() {
		var txnID int64
		err = rows.Scan(&txnID)
		if err != nil {
			log.Errorf("rows.Scan: %v", err)
		}
		txnIDs = append(txnIDs, txnID)
	}

	json.NewEncoder(w).Encode(txnIDs)
}

func GetTransactionSum(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var txnID string
	var ok bool
	if txnID, ok = params["transaction_id"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("transaction_id not present"))
		return
	}

	row := db.QueryRowContext(r.Context(), sqlGetSum, txnID)
	if row.Err() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("GetTransactionSum: db.QueryRowContext: %v", row.Err())
		return
	}

	var output int64
	err := row.Scan(&output)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		log.Infof("GetTransactionSum: row.Scan sql.ErrNoRows: %v", row.Err())
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("GetTransactionSum: row.Scan: %v", err)
		return
	}

	json.NewEncoder(w).Encode(output)
}
