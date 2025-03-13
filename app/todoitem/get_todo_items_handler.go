package todoitem

import (
	"context"

	"github.com/nsaltun/todolist-service/domain"
)

type GetTodoItemsRequest struct {
	ID string `json:"id"`
}

type GetTodoItemsResponse struct {
	items       []domain.TodoItem
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetTodoItemsHandler struct {
	repo Repository
}

func NewGetTodoItemsHandler(repo Repository) *GetTodoItemsHandler {
	return &GetTodoItemsHandler{repo: repo}
}

func (h *GetTodoItemsHandler) Handle(ctx context.Context, r *GetTodoItemsRequest) (*GetTodoItemsResponse, error) {
	todoItems, err := h.repo.GetTodoItems(ctx)
	if err != nil {
		return nil, err
	}

	return &GetTodoItemsResponse{
		items: todoItems,
	}, nil
}
