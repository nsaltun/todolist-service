package todoitem

import (
	"context"

	"github.com/nsaltun/todolist-service/domain"
	"github.com/nsaltun/todolist-service/pkg/pagination"
)

type GetTodoItemsRequest struct {
	Status     *string `json:"status,omitempty"`
	SearchTerm *string `json:"search_term,omitempty"`
	pagination.PaginationRequest
}

type GetTodoItemsResponse struct {
	Items      []domain.TodoItem
	Pagination pagination.PaginationResponse `json:"pagination"`
}

type GetTodoItemsHandler struct {
	repo Repository
}

func NewGetTodoItemsHandler(repo Repository) *GetTodoItemsHandler {
	return &GetTodoItemsHandler{repo: repo}
}

func (h *GetTodoItemsHandler) Handle(ctx context.Context, r *GetTodoItemsRequest) (*GetTodoItemsResponse, error) {
	filter := domain.TodoFilter{
		Status:     r.Status,
		SearchTerm: r.SearchTerm,
		Pagination: pagination.Pagination{
			Limit:  r.Limit,
			Offset: r.Offset,
		},
	}

	items, total, err := h.repo.GetTodoItems(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &GetTodoItemsResponse{
		Items:      items,
		Pagination: pagination.NewPaginationResponse(r.Limit, r.Offset, total),
	}, nil
}
