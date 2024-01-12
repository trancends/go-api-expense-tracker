package config

const (
	SelectLastInsert = "SELECT balance FROM expenses ORDER BY created_at DESC LIMIT 1"
	InsertExpense    = `INSERT INTO expenses 
    (date, amount, transaction_type, balance, descritption,created_at,updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
  `
	SelectExpenseBetwenDate = `
    SELECT * FROM expenses 
    WHERE date >= $1 AND date <= $2 LIMIT $3 OFFSET $4 `
	SelectExpenseByID   = `SELECT * FROM expenses WHERE id = $1`
	SelectExpenseByType = `
    SELECT * FROM expenses WHERE 
    transaction_type = $1`
)
