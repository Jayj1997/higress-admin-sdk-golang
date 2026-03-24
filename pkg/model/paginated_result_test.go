// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPaginatedResult(t *testing.T) {
	tests := []struct {
		name        string
		data        []string
		total       int
		pageNum     int
		pageSize    int
		expectData  int
		expectPages int
	}{
		{
			name:        "single page",
			data:        []string{"a", "b", "c"},
			total:       3,
			pageNum:     1,
			pageSize:    10,
			expectData:  3,
			expectPages: 1,
		},
		{
			name:        "multiple pages - first page",
			data:        []string{"a", "b"},
			total:       5,
			pageNum:     1,
			pageSize:    2,
			expectData:  2,
			expectPages: 3,
		},
		{
			name:        "multiple pages - last page",
			data:        []string{"e"},
			total:       5,
			pageNum:     3,
			pageSize:    2,
			expectData:  1,
			expectPages: 3,
		},
		{
			name:        "empty data",
			data:        []string{},
			total:       0,
			pageNum:     1,
			pageSize:    10,
			expectData:  0,
			expectPages: 0,
		},
		{
			name:        "exact page boundary",
			data:        []string{"a", "b"},
			total:       4,
			pageNum:     2,
			pageSize:    2,
			expectData:  2,
			expectPages: 2,
		},
		{
			name:        "zero page size",
			data:        []string{"a", "b", "c"},
			total:       3,
			pageNum:     1,
			pageSize:    0,
			expectData:  3,
			expectPages: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPaginatedResult(tt.data, tt.total, tt.pageNum, tt.pageSize)

			require.NotNil(t, result)
			assert.Len(t, result.Data, tt.expectData)
			assert.Equal(t, tt.total, result.Total)
			assert.Equal(t, tt.pageNum, result.PageNum)
			assert.Equal(t, tt.pageSize, result.PageSize)
			assert.Equal(t, tt.expectPages, result.TotalPages)
		})
	}
}

func TestPaginatedResultFromFullList(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	tests := []struct {
		name           string
		list           []string
		query          *CommonPageQuery
		expectData     int
		expectTotal    int
		expectPageNum  int
		expectPageSize int
	}{
		{
			name:           "nil query returns all",
			list:           items,
			query:          nil,
			expectData:     10,
			expectTotal:    10,
			expectPageNum:  1,
			expectPageSize: 10,
		},
		{
			name:           "first page",
			list:           items,
			query:          &CommonPageQuery{PageNum: 1, PageSize: 3},
			expectData:     3,
			expectTotal:    10,
			expectPageNum:  1,
			expectPageSize: 3,
		},
		{
			name:           "middle page",
			list:           items,
			query:          &CommonPageQuery{PageNum: 2, PageSize: 3},
			expectData:     3,
			expectTotal:    10,
			expectPageNum:  2,
			expectPageSize: 3,
		},
		{
			name:           "last page partial",
			list:           items,
			query:          &CommonPageQuery{PageNum: 4, PageSize: 3},
			expectData:     1,
			expectTotal:    10,
			expectPageNum:  4,
			expectPageSize: 3,
		},
		{
			name:           "page beyond data",
			list:           items,
			query:          &CommonPageQuery{PageNum: 10, PageSize: 3},
			expectData:     0,
			expectTotal:    10,
			expectPageNum:  10,
			expectPageSize: 3,
		},
		{
			name:           "empty list",
			list:           []string{},
			query:          &CommonPageQuery{PageNum: 1, PageSize: 3},
			expectData:     0,
			expectTotal:    0,
			expectPageNum:  1,
			expectPageSize: 3,
		},
		{
			name:           "zero page num defaults to 1",
			list:           items,
			query:          &CommonPageQuery{PageNum: 0, PageSize: 5},
			expectData:     5,
			expectTotal:    10,
			expectPageNum:  1,
			expectPageSize: 5,
		},
		{
			name:           "negative page num defaults to 1",
			list:           items,
			query:          &CommonPageQuery{PageNum: -1, PageSize: 5},
			expectData:     5,
			expectTotal:    10,
			expectPageNum:  1,
			expectPageSize: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PaginatedResultFromFullList(tt.list, tt.query)

			require.NotNil(t, result)
			assert.Len(t, result.Data, tt.expectData)
			assert.Equal(t, tt.expectTotal, result.Total)
			assert.Equal(t, tt.expectPageNum, result.PageNum)
			assert.Equal(t, tt.expectPageSize, result.PageSize)
		})
	}
}

func TestCommonPageQuery_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		query    *CommonPageQuery
		expected int
	}{
		{"page 1, size 10", &CommonPageQuery{PageNum: 1, PageSize: 10}, 0},
		{"page 2, size 10", &CommonPageQuery{PageNum: 2, PageSize: 10}, 10},
		{"page 3, size 20", &CommonPageQuery{PageNum: 3, PageSize: 20}, 40},
		{"page 0 defaults to offset 0", &CommonPageQuery{PageNum: 0, PageSize: 10}, 0},
		{"negative page defaults to offset 0", &CommonPageQuery{PageNum: -1, PageSize: 10}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.query.GetOffset())
		})
	}
}

func TestCommonPageQuery_GetPageSize(t *testing.T) {
	tests := []struct {
		name     string
		query    *CommonPageQuery
		expected int
	}{
		{"valid page size", &CommonPageQuery{PageSize: 20}, 20},
		{"zero page size defaults to 10", &CommonPageQuery{PageSize: 0}, 10},
		{"negative page size defaults to 10", &CommonPageQuery{PageSize: -5}, 10},
		{"page size over 100 capped to 100", &CommonPageQuery{PageSize: 200}, 100},
		{"page size exactly 100", &CommonPageQuery{PageSize: 100}, 100},
		{"page size exactly 1", &CommonPageQuery{PageSize: 1}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.query.GetPageSize())
		})
	}
}

func TestRoutePageQuery(t *testing.T) {
	query := &RoutePageQuery{
		CommonPageQuery: CommonPageQuery{PageNum: 2, PageSize: 20},
		DomainName:      "example.com",
		Name:            "test-route",
	}

	assert.Equal(t, 2, query.PageNum)
	assert.Equal(t, 20, query.PageSize)
	assert.Equal(t, "example.com", query.DomainName)
	assert.Equal(t, "test-route", query.Name)
	assert.Equal(t, 20, query.GetPageSize())
	assert.Equal(t, 20, query.GetOffset())
}

func TestWasmPluginPageQuery(t *testing.T) {
	builtIn := true
	query := &WasmPluginPageQuery{
		CommonPageQuery: CommonPageQuery{PageNum: 1, PageSize: 10},
		Name:            "basic-auth",
		Version:         "1.0.0",
		Category:        "security",
		BuiltIn:         &builtIn,
		Lang:            "zh-CN",
	}

	assert.Equal(t, 1, query.PageNum)
	assert.Equal(t, 10, query.PageSize)
	assert.Equal(t, "basic-auth", query.Name)
	assert.Equal(t, "1.0.0", query.Version)
	assert.Equal(t, "security", query.Category)
	require.NotNil(t, query.BuiltIn)
	assert.True(t, *query.BuiltIn)
	assert.Equal(t, "zh-CN", query.Lang)
}

func TestWasmPluginPageQuery_BuiltInFalse(t *testing.T) {
	builtIn := false
	query := &WasmPluginPageQuery{
		BuiltIn: &builtIn,
	}

	require.NotNil(t, query.BuiltIn)
	assert.False(t, *query.BuiltIn)
}

func TestWasmPluginPageQuery_NilBuiltIn(t *testing.T) {
	query := &WasmPluginPageQuery{}

	assert.Nil(t, query.BuiltIn)
}

func TestPaginatedResult_GenericType(t *testing.T) {
	t.Run("with int type", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		result := NewPaginatedResult(data, 5, 1, 10)

		require.NotNil(t, result)
		assert.Len(t, result.Data, 5)
		assert.Equal(t, 5, result.Total)
	})

	t.Run("with struct type", func(t *testing.T) {
		type TestItem struct {
			ID   string
			Name string
		}
		data := []TestItem{
			{ID: "1", Name: "Item 1"},
			{ID: "2", Name: "Item 2"},
		}
		result := NewPaginatedResult(data, 2, 1, 10)

		require.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, "1", result.Data[0].ID)
		assert.Equal(t, "Item 2", result.Data[1].Name)
	})
}

func TestPaginatedResultFromFullList_PaginationCorrectness(t *testing.T) {
	items := make([]int, 100)
	for i := 0; i < 100; i++ {
		items[i] = i + 1
	}

	// Test first page
	result := PaginatedResultFromFullList(items, &CommonPageQuery{PageNum: 1, PageSize: 10})
	require.Len(t, result.Data, 10)
	assert.Equal(t, 1, result.Data[0])
	assert.Equal(t, 10, result.Data[9])

	// Test middle page
	result = PaginatedResultFromFullList(items, &CommonPageQuery{PageNum: 5, PageSize: 10})
	require.Len(t, result.Data, 10)
	assert.Equal(t, 41, result.Data[0])
	assert.Equal(t, 50, result.Data[9])

	// Test last page
	result = PaginatedResultFromFullList(items, &CommonPageQuery{PageNum: 10, PageSize: 10})
	require.Len(t, result.Data, 10)
	assert.Equal(t, 91, result.Data[0])
	assert.Equal(t, 100, result.Data[9])
}
