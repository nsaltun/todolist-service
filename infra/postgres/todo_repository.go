package postgres

import (
	"context"

	"github.com/nsaltun/todolist-service/domain"
)

type TodoRepository struct {
	conn *PostgresConnection
}

func NewTodoRepository(conn *PostgresConnection) *TodoRepository {
	return &TodoRepository{conn: conn}
}

func (r *TodoRepository) Create(ctx context.Context, todo domain.TodoItem) error {
	query := `INSERT INTO todos (id, title, description, status) VALUES ($1, $2, $3, $4)`
	_, err := r.conn.dbPool.Exec(ctx, query, todo.ID, todo.Title, todo.Description, todo.Status)
	return err
}

func (r *TodoRepository) GetTodoItems(ctx context.Context) ([]domain.TodoItem, error) {
	r.conn.dbPool.Query(ctx, "SELECT * FROM todo_items")
	return []domain.TodoItem{}, nil
}
