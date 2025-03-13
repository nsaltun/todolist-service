package todoitem

import (
	"context"

	"github.com/nsaltun/todolist-service/domain"
)

type Repository interface {
	GetTodoItems(ctx context.Context) ([]domain.TodoItem, error)
	Create(ctx context.Context, todo domain.TodoItem) error
}
