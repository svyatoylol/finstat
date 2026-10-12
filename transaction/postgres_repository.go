package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var _ Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

const selectColumns = `id, user_id, type, amount, category`

func scanTransaction(row interface{ Scan(...any) error }) (Transaction, error) {
	var t Transaction
	err := row.Scan(&t.ID, &t.UserID, &t.Type, &t.Amount, &t.Category)
	return t, err
}

func (r *PostgresRepository) Add(ctx context.Context, t Transaction) (Transaction, error) {
	const q = `INSERT INTO transactions (user_id, type, amount, category)
	           VALUES ($1, $2, $3, $4)
	           RETURNING ` + selectColumns

	out, err := scanTransaction(r.db.QueryRowContext(ctx, q, t.UserID, t.Type, t.Amount, t.Category))
	if err != nil {
		return Transaction{}, fmt.Errorf("insert transaction: %w", err)
	}
	return out, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, userID, id int64) (Transaction, error) {
	const q = `SELECT ` + selectColumns + ` FROM transactions WHERE id = $1 AND user_id = $2`

	t, err := scanTransaction(r.db.QueryRowContext(ctx, q, id, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return Transaction{}, ErrNotFound
	}
	if err != nil {
		return Transaction{}, fmt.Errorf("get transaction %d: %w", id, err)
	}
	return t, nil
}

func (r *PostgresRepository) GetFiltered(ctx context.Context, f Filter) ([]Transaction, error) {
	query := `SELECT ` + selectColumns + ` FROM transactions WHERE user_id = $1`
	args := []any{f.UserID}

	if f.Type != nil {
		args = append(args, *f.Type)
		query += fmt.Sprintf(" AND type = $%d", len(args))
	}
	if f.Category != nil {
		args = append(args, *f.Category)
		query += fmt.Sprintf(" AND category = $%d", len(args))
	}
	if f.AfterID > 0 {
		args = append(args, f.AfterID)
		query += fmt.Sprintf(" AND id > $%d", len(args))
	}
	query += " ORDER BY id"
	limit := f.Limit
	if limit <= 0 {
		limit = DefaultPageSize
	}
	args = append(args, limit)
	query += fmt.Sprintf(" LIMIT $%d", len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var results []Transaction
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		results = append(results, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transactions: %w", err)
	}
	return results, nil
}

func (r *PostgresRepository) GetBalance(ctx context.Context, userID int64) (float64, error) {
	const q = `SELECT COALESCE(SUM(CASE
		WHEN type = 'incoming' THEN amount
		WHEN type = 'outgoing' THEN -amount
		ELSE 0
	END), 0) FROM transactions WHERE user_id = $1`

	var balance float64
	if err := r.db.QueryRowContext(ctx, q, userID).Scan(&balance); err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}
	return balance, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, userID, id int64) error {
	const q = `DELETE FROM transactions WHERE id = $1 AND user_id = $2`

	res, err := r.db.ExecContext(ctx, q, id, userID)
	if err != nil {
		return fmt.Errorf("delete transaction %d: %w", id, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, userID int64, t Transaction) (Transaction, error) {
	const q = `UPDATE transactions SET type = $1, amount = $2, category = $3
	           WHERE id = $4 AND user_id = $5
	           RETURNING ` + selectColumns

	out, err := scanTransaction(r.db.QueryRowContext(ctx, q, t.Type, t.Amount, t.Category, t.ID, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return Transaction{}, ErrNotFound
	}
	if err != nil {
		return Transaction{}, fmt.Errorf("update transaction %d: %w", t.ID, err)
	}
	return out, nil
}
