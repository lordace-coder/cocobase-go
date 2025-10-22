package tests

import (
	"testing"

	"github.com/yourusername/cocobase-go/cocobase"
)

func TestQueryBuilder(t *testing.T) {
	tests := []struct {
		name     string
		builder  func() *cocobase.QueryBuilder
		expected string
	}{
		{
			name: "simple where",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().Where("status", "active")
			},
			expected: "status=active",
		},
		{
			name: "with operator",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().Filter("age", "gte", 18)
			},
			expected: "age_gte=18",
		},
		{
			name: "simple OR",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().
					Or("isPremium", "eq", true).
					Or("isVerified", "eq", true)
			},
			expected: "%5Bor%5DisPremium=true&%5Bor%5DisVerified=true",
		},
		{
			name: "with pagination",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().
					Where("status", "active").
					Limit(50).
					Offset(100)
			},
			expected: "limit=50&offset=100&status=active",
		},
		{
			name: "with sorting",
			builder: func() *cocobase.QueryBuilder {
				return cocobase.NewQuery().
					Sort("createdAt").
					OrderDesc()
			},
			expected: "order=desc&sort=createdAt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder().Build()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMultiFieldOr(t *testing.T) {
	query := cocobase.NewQuery().
		MultiFieldOr([]string{"name", "email"}, "contains", "john")

	result := query.Build()
	expected := "name__or__email_contains=john"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestOrGroups(t *testing.T) {
	query := cocobase.NewQuery().
		OrGroup("tier", "isPremium", "eq", true).
		OrGroup("tier", "isVerified", "eq", true).
		OrGroup("location", "country", "eq", "US").
		OrGroup("location", "country", "eq", "UK")

	result := query.Build()
	
	// Check if all OR group parameters are present
	if !contains(result, "%5Bor%3Atier%5D") {
		t.Error("Missing tier OR group")
	}
	if !contains(result, "%5Bor%3Alocation%5D") {
		t.Error("Missing location OR group")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
