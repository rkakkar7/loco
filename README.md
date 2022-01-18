### Table Schema

```
CREATE TABLE transaction (
	transaction_id serial NOT NULL,
	amount decimal NOT NULL DEFAULT 0,
	type VARCHAR(20),
	parent_id serial,
	PRIMARY KEY (transaction_id),
	FOREIGN KEY (parent_id) REFERENCES transaction (transaction_id)
);
```

### Routes

- PUT "/transactionservice/transaction/$transaction_id"
- GET "/transactionservice/transaction/$transaction_id"
- GET "/transactionservice/types/$type"
- GET "/transactionservice/sum/$transaction_id"


### Installation Steps
1. Install Postgres
2. Create Table using above mentioned schema
3. Run using `go run main.go` or build binary using `go build main.go` and spin up a background process
4. Server starts at `127.0.0.1:8080` 