package postgres

import (
	"context"
	"fmt"

	"github.com/nsaltun/todolist-service/domain"
)

type TodoRepository struct {
	conn *PostgresConnection
}

func NewTodoRepository(conn *PostgresConnection) *TodoRepository {
	return &TodoRepository{conn: conn}
}

func (r *TodoRepository) Create(ctx context.Context, todo domain.TodoItem) error {
	query := `INSERT INTO todo_items (id, title, description, status, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := r.conn.dbPool.Exec(ctx, query, todo.ID, todo.Title, todo.Description, todo.Status)
	return err
}

func (r *TodoRepository) GetTodoItems(ctx context.Context, filter domain.TodoFilter) ([]domain.TodoItem, int64, error) {
	// First, get total count
	countQuery := `
        SELECT COUNT(*) 
        FROM todo_items 
        WHERE 1=1`

	var args []interface{}
	if filter.Status != nil {
		args = append(args, *filter.Status)
		countQuery += fmt.Sprintf(" AND status = $%d", len(args))
	}

	if filter.SearchTerm != nil {
		args = append(args, "%"+*filter.SearchTerm+"%") // Add wildcards for ILIKE
		countQuery += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)",
			len(args), len(args))
	}

	var total int64
	err := r.conn.dbPool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Then get paginated items
	query := `
        SELECT id, title, description, status, created_at, updated_at
        FROM todo_items
        WHERE 1=1`

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", len(args)-1) // Reuse previous args
	}

	if filter.SearchTerm != nil {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)",
			len(args), len(args))
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	// Add pagination parameters to args
	args = append(args, filter.Pagination.GetLimit(), filter.Pagination.GetOffset())

	rows, err := r.conn.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query todo items: %w", err)
	}
	defer rows.Close()

	var todoItems []domain.TodoItem
	for rows.Next() {
		var todo domain.TodoItem
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Status,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan todo item: %w", err)
		}
		todoItems = append(todoItems, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating todo items: %w", err)
	}

	return todoItems, total, nil
}
