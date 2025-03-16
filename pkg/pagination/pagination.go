package pagination

// Pagination represents the core pagination parameters
type Pagination struct {
	Limit  int
	Offset int
}

func (p Pagination) GetLimit() int {
	if p.Limit <= 0 {
		return 10 // default limit
	}
	if p.Limit > 100 {
		return 100 // max limit
	}
	return p.Limit
}

func (p Pagination) GetOffset() int {
	if p.Offset < 0 {
		return 0
	}
	return p.Offset
}

// PaginationRequest represents the API request pagination parameters
type PaginationRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// PaginationResponse represents the API response pagination metadata
type PaginationResponse struct {
	Limit     int   `json:"limit"`
	Offset    int   `json:"offset"`
	Total     int64 `json:"total"`
	HasNext   bool  `json:"has_next"`
	HasBefore bool  `json:"has_before"`
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(limit int, offset int, total int64) PaginationResponse {
	return PaginationResponse{
		Limit:     limit,
		Offset:    offset,
		Total:     total,
		HasNext:   int64(offset+limit) < total,
		HasBefore: offset > 0,
	}
}
