// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRouteService_Interface tests that the implementation satisfies the interface
func TestRouteService_Interface(t *testing.T) {
	var _ RouteService = (*RouteServiceImpl)(nil)
}

// TestRouteService_New tests creating a new RouteService
func TestRouteService_New(t *testing.T) {
	// Create service with nil dependencies (for testing without cluster)
	svc := NewRouteService(nil, nil, nil)
	require.NotNil(t, svc)
}

// TestRouteModel tests the Route model
func TestRouteModel(t *testing.T) {
	pathPredicate := &route.RoutePredicate{
		MatchType: route.MatchTypePrefix,
		Path:      "/api/v1",
	}
	r := model.Route{
		Name:    "test-route",
		Path:    pathPredicate,
		Domains: []string{"example.com", "api.example.com"},
	}

	assert.Equal(t, "test-route", r.Name)
	assert.Equal(t, route.MatchTypePrefix, r.Path.MatchType)
	assert.Equal(t, "/api/v1", r.Path.Path)
	assert.Len(t, r.Domains, 2)
	assert.Contains(t, r.Domains, "example.com")
}

// TestRouteModel_Validation tests route validation
func TestRouteModel_Validation(t *testing.T) {
	tests := []struct {
		name        string
		route       *model.Route
		expectError bool
	}{
		{
			name: "valid route with path and services",
			route: &model.Route{
				Name: "valid-route",
				Path: &route.RoutePredicate{
					MatchType: route.MatchTypePrefix,
					Path:      "/api",
				},
				Services: []*route.UpstreamService{
					{Name: "backend-service"},
				},
			},
			expectError: false,
		},
		{
			name: "route with domains",
			route: &model.Route{
				Name: "route-with-domains",
				Path: &route.RoutePredicate{
					MatchType: route.MatchTypePrefix,
					Path:      "/api",
				},
				Domains: []string{"example.com"},
				Services: []*route.UpstreamService{
					{Name: "backend-service"},
				},
			},
			expectError: false,
		},
		{
			name: "route with empty name",
			route: &model.Route{
				Name: "",
				Path: &route.RoutePredicate{
					MatchType: route.MatchTypePrefix,
					Path:      "/api",
				},
				Services: []*route.UpstreamService{
					{Name: "backend-service"},
				},
			},
			expectError: true,
		},
		{
			name: "route without services",
			route: &model.Route{
				Name: "route-no-services",
				Path: &route.RoutePredicate{
					MatchType: route.MatchTypePrefix,
					Path:      "/api",
				},
				Services: nil,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.route.Validate()
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// TestRoutePageQuery tests the RoutePageQuery model
func TestRoutePageQuery(t *testing.T) {
	query := &model.RoutePageQuery{
		CommonPageQuery: model.CommonPageQuery{
			PageNum:  1,
			PageSize: 10,
		},
		DomainName: "example.com",
	}

	assert.Equal(t, 1, query.PageNum)
	assert.Equal(t, 10, query.PageSize)
	assert.Equal(t, "example.com", query.DomainName)
}

// TestRoutePagination tests pagination logic for routes
func TestRoutePagination(t *testing.T) {
	routes := make([]model.Route, 25)
	for i := 0; i < 25; i++ {
		routes[i] = model.Route{
			Name: "route-" + string(rune('a'+i)),
			Path: &route.RoutePredicate{
				MatchType: route.MatchTypePrefix,
				Path:      "/path-" + string(rune('a'+i)),
			},
		}
	}

	// Test first page
	total := len(routes)
	pageNum := 1
	pageSize := 10
	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pagedData := routes[start:end]
	result := model.NewPaginatedResult(pagedData, total, pageNum, pageSize)

	assert.Equal(t, 10, len(result.Data))
	assert.Equal(t, 25, result.Total)
	assert.Equal(t, 1, result.PageNum)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, 3, result.TotalPages)
}

// TestRouteWithDomains tests route with multiple domains
func TestRouteWithDomains(t *testing.T) {
	r := model.Route{
		Name:    "multi-domain-route",
		Path:    &route.RoutePredicate{MatchType: route.MatchTypePrefix, Path: "/api"},
		Domains: []string{"api1.example.com", "api2.example.com", "api3.example.com"},
	}

	assert.Len(t, r.Domains, 3)
	for _, domain := range r.Domains {
		assert.Contains(t, domain, "example.com")
	}
}

// TestRouteWithServices tests route with services
func TestRouteWithServices(t *testing.T) {
	weight80 := 80
	weight20 := 20
	r := model.Route{
		Name: "route-with-services",
		Path: &route.RoutePredicate{MatchType: route.MatchTypePrefix, Path: "/api"},
		Services: []*route.UpstreamService{
			{Name: "service-a", Weight: &weight80},
			{Name: "service-b", Weight: &weight20},
		},
	}

	assert.Len(t, r.Services, 2)
	assert.Equal(t, "service-a", r.Services[0].Name)
	assert.Equal(t, 80, *r.Services[0].Weight)
}

// TestUpstreamService tests the UpstreamService model
func TestUpstreamService(t *testing.T) {
	weight := 100
	port := 8080
	svc := &route.UpstreamService{
		Name:   "backend-service",
		Weight: &weight,
		Port:   port,
	}

	assert.Equal(t, "backend-service", svc.Name)
	assert.Equal(t, 100, *svc.Weight)
	assert.Equal(t, 8080, svc.Port)
}

// TestRoutePredicate tests the RoutePredicate model
func TestRoutePredicate(t *testing.T) {
	tests := []struct {
		name        string
		predicate   *route.RoutePredicate
		matchType   string
		path        string
		expectValid bool
	}{
		{
			name:        "prefix path",
			predicate:   &route.RoutePredicate{MatchType: route.MatchTypePrefix, Path: "/api"},
			matchType:   route.MatchTypePrefix,
			path:        "/api",
			expectValid: true,
		},
		{
			name:        "exact path",
			predicate:   &route.RoutePredicate{MatchType: route.MatchTypeExact, Path: "/health"},
			matchType:   route.MatchTypeExact,
			path:        "/health",
			expectValid: true,
		},
		{
			name:        "regex path",
			predicate:   &route.RoutePredicate{MatchType: route.MatchTypeRegex, Path: "^/api/v[0-9]+/.*"},
			matchType:   route.MatchTypeRegex,
			path:        "^/api/v[0-9]+/.*",
			expectValid: true,
		},
		{
			name:        "empty matchType",
			predicate:   &route.RoutePredicate{MatchType: "", Path: "/api"},
			matchType:   "",
			path:        "/api",
			expectValid: false,
		},
		{
			name:        "empty path",
			predicate:   &route.RoutePredicate{MatchType: route.MatchTypePrefix, Path: ""},
			matchType:   route.MatchTypePrefix,
			path:        "",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.matchType, tt.predicate.MatchType)
			assert.Equal(t, tt.path, tt.predicate.Path)

			err := tt.predicate.Validate()
			if tt.expectValid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

// TestMatchTypeConstants tests the match type constants
func TestMatchTypeConstants(t *testing.T) {
	assert.Equal(t, "exact", route.MatchTypeExact)
	assert.Equal(t, "prefix", route.MatchTypePrefix)
	assert.Equal(t, "regex", route.MatchTypeRegex)
}
