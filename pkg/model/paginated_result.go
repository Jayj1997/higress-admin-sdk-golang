// Package model provides data models for Higress Admin SDK.
package model

// PaginatedResult represents a paginated result set.
type PaginatedResult[T any] struct {
	// Data is the list of items in the current page.
	Data []T `json:"data,omitempty"`

	// Total is the total number of items.
	Total int `json:"total,omitempty"`

	// PageNum is the current page number (1-based).
	PageNum int `json:"pageNum,omitempty"`

	// PageSize is the number of items per page.
	PageSize int `json:"pageSize,omitempty"`

	// TotalPages is the total number of pages.
	TotalPages int `json:"totalPages,omitempty"`
}

// NewPaginatedResult creates a new PaginatedResult.
func NewPaginatedResult[T any](data []T, total, pageNum, pageSize int) *PaginatedResult[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	return &PaginatedResult[T]{
		Data:       data,
		Total:      total,
		PageNum:    pageNum,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// CommonPageQuery represents common pagination query parameters.
type CommonPageQuery struct {
	// PageNum is the page number (1-based).
	PageNum int `json:"pageNum,omitempty"`

	// PageSize is the number of items per page.
	PageSize int `json:"pageSize,omitempty"`
}

// GetOffset returns the offset for database queries.
func (q *CommonPageQuery) GetOffset() int {
	if q.PageNum <= 0 {
		return 0
	}
	return (q.PageNum - 1) * q.GetPageSize()
}

// GetPageSize returns the page size with a default value.
func (q *CommonPageQuery) GetPageSize() int {
	if q.PageSize <= 0 {
		return 10
	}
	if q.PageSize > 100 {
		return 100
	}
	return q.PageSize
}

// RoutePageQuery represents route query parameters.
type RoutePageQuery struct {
	CommonPageQuery

	// DomainName filters routes by domain name.
	DomainName string `json:"domainName,omitempty"`

	// Name filters routes by name.
	Name string `json:"name,omitempty"`
}

// WasmPluginPageQuery represents WASM plugin query parameters.
type WasmPluginPageQuery struct {
	CommonPageQuery

	// Name filters plugins by name.
	Name string `json:"name,omitempty"`

	// Category filters plugins by category.
	Category string `json:"category,omitempty"`

	// BuiltIn filters plugins by built-in status.
	BuiltIn *bool `json:"builtIn,omitempty"`
}
