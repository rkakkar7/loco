package api

type Transaction struct {
	Amount   float64 `json:"amount" field:"amount"`
	Type     string  `json:"type" field:"type"`
	ParentID *int64  `json:"parent_id,omitempty" field:"parent_id"`
}
