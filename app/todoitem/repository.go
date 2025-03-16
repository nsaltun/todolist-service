package todoitem

import (
	"context"

	"github.com/nsaltun/todolist-service/domain"
)

type Repository interface {
	GetTodoItems(ctx context.Context, filter domain.TodoFilter) ([]domain.TodoItem, int64, error)
	Create(ctx context.Context, todo domain.TodoItem) error
}
