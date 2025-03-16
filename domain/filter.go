package domain

import "github.com/nsaltun/todolist-service/pkg/pagination"

type TodoFilter struct {
	Status     *string
	SearchTerm *string
	pagination.Pagination
}
