package todoitem

import (
	"context"

	"github.com/google/uuid"
	"github.com/nsaltun/todolist-service/domain"
)

type CreateTodoItemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTodoItemResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type CreateTodoItemHandler struct {
	repo Repository
}

func NewCreateTodoItemHandler(repo Repository) *CreateTodoItemHandler {
	return &CreateTodoItemHandler{repo: repo}
}

func (h *CreateTodoItemHandler) Handle(ctx context.Context, r *CreateTodoItemRequest) (*CreateTodoItemResponse, error) {
	todo := domain.TodoItem{
		ID:          uuid.New().String(),
		Title:       r.Title,
		Description: r.Description,
		Status:      "pending",
	}

	err := h.repo.Create(ctx, todo)
	if err != nil {
		return nil, err
	}

	return &CreateTodoItemResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
	}, nil
}
