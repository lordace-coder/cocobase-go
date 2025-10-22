package tests

import (
	"net/url"
	"strings"
	"testing"

	"github.com/lordace-coder/cocobase-go/cocobase"
)

// Helper function to parse query string
func parseQuery(queryStr string) url.Values {
	params, _ := url.ParseQuery(queryStr)
	return params
}

// Helper to check if query contains parameter
func hasParam(query, key, value string) bool {
	params := parseQuery(query)
	return params.Get(key) == value
}

// ============================================
// 1. BASIC COMPARISON OPERATORS
// ============================================

func TestWhere(t *testing.T) {
	query := cocobase.NewQuery().
		Where("status", "active")

	result := query.Build()

	if !hasParam(result, "status", "active") {
		t.Errorf("Expected status=active, got %s", result)
	}
}

func TestEquals(t *testing.T) {
	query := cocobase.NewQuery().
		Equals("role", "admin")

	result := query.Build()

	if !hasParam(result, "role", "admin") {
		t.Errorf("Expected role=admin, got %s", result)
	}
}

func TestNotEquals(t *testing.T) {
	query := cocobase.NewQuery().
		NotEquals("status", "banned")

	result := query.Build()

	if !hasParam(result, "status_ne", "banned") {
		t.Errorf("Expected status_ne=banned, got %s", result)
	}
}

func TestGreaterThan(t *testing.T) {
	query := cocobase.NewQuery().
		GreaterThan("age", 18)

	result := query.Build()

	if !hasParam(result, "age_gt", "18") {
		t.Errorf("Expected age_gt=18, got %s", result)
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	query := cocobase.NewQuery().
		GreaterThanOrEqual("age", 18)

	result := query.Build()

	if !hasParam(result, "age_gte", "18") {
		t.Errorf("Expected age_gte=18, got %s", result)
	}
}

func TestLessThan(t *testing.T) {
	query := cocobase.NewQuery().
		LessThan("price", 100)

	result := query.Build()

	if !hasParam(result, "price_lt", "100") {
		t.Errorf("Expected price_lt=100, got %s", result)
	}
}

func TestLessThanOrEqual(t *testing.T) {
	query := cocobase.NewQuery().
		LessThanOrEqual("stock", 10)

	result := query.Build()

	if !hasParam(result, "stock_lte", "10") {
		t.Errorf("Expected stock_lte=10, got %s", result)
	}
}

func TestBetween(t *testing.T) {
	query := cocobase.NewQuery().
		Between("age", 18, 65)

	result := query.Build()
	params := parseQuery(result)

	if params.Get("age_gte") != "18" {
		t.Errorf("Expected age_gte=18")
	}
	if params.Get("age_lte") != "65" {
		t.Errorf("Expected age_lte=65")
	}
}

// ============================================
// 2. STRING OPERATORS
// ============================================

func TestContains(t *testing.T) {
	query := cocobase.NewQuery().
		Contains("name", "john")

	result := query.Build()

	if !hasParam(result, "name_contains", "john") {
		t.Errorf("Expected name_contains=john, got %s", result)
	}
}

func TestStartsWith(t *testing.T) {
	query := cocobase.NewQuery().
		StartsWith("email", "admin")

	result := query.Build()

	if !hasParam(result, "email_startswith", "admin") {
		t.Errorf("Expected email_startswith=admin, got %s", result)
	}
}

func TestEndsWith(t *testing.T) {
	query := cocobase.NewQuery().
		EndsWith("email", "gmail.com")

	result := query.Build()

	if !hasParam(result, "email_endswith", "gmail.com") {
		t.Errorf("Expected email_endswith=gmail.com, got %s", result)
	}
}

func TestSearch(t *testing.T) {
	query := cocobase.NewQuery().
		Search("admin", "name", "email", "username")

	result := query.Build()

	if !strings.Contains(result, "name__or__email__or__username_contains=admin") {
		t.Errorf("Expected multi-field search, got %s", result)
	}
}

// ============================================
// 3. LIST OPERATORS
// ============================================

func TestIn(t *testing.T) {
	query := cocobase.NewQuery().
		In("role", "admin", "moderator", "support")

	result := query.Build()

	if !hasParam(result, "role_in", "admin,moderator,support") {
		t.Errorf("Expected role_in=admin,moderator,support, got %s", result)
	}
}

func TestNotIn(t *testing.T) {
	query := cocobase.NewQuery().
		NotIn("status", "banned", "deleted")

	result := query.Build()

	if !hasParam(result, "status_notin", "banned,deleted") {
		t.Errorf("Expected status_notin=banned,deleted, got %s", result)
	}
}

func TestIsNull(t *testing.T) {
	query := cocobase.NewQuery().
		IsNull("deletedAt")

	result := query.Build()

	if !hasParam(result, "deletedAt_isnull", "true") {
		t.Errorf("Expected deletedAt_isnull=true, got %s", result)
	}
}

func TestIsNotNull(t *testing.T) {
	query := cocobase.NewQuery().
		IsNotNull("profilePicture")

	result := query.Build()

	if !hasParam(result, "profilePicture_isnull", "false") {
		t.Errorf("Expected profilePicture_isnull=false, got %s", result)
	}
}

// ============================================
// 4. SIMPLE OR CONDITIONS
// ============================================

func TestSimpleOr(t *testing.T) {
	query := cocobase.NewQuery().
		Or().
		Where("isPremium", true).
		Where("isVerified", true).
		Done()

	result := query.Build()

	if !strings.Contains(result, "%5Bor%5DisPremium=true") {
		t.Errorf("Expected [or]isPremium=true in query, got %s", result)
	}
	if !strings.Contains(result, "%5Bor%5DisVerified=true") {
		t.Errorf("Expected [or]isVerified=true in query, got %s", result)
	}
}

func TestOrWithOperators(t *testing.T) {
	query := cocobase.NewQuery().
		Or().
		GreaterThan("age", 65).
		Equals("role", "admin").
		Done()

	result := query.Build()

	if !strings.Contains(result, "%5Bor%5Dage_gt=65") {
		t.Errorf("Expected [or]age_gt=65 in query, got %s", result)
	}
	if !strings.Contains(result, "%5Bor%5Drole=admin") {
		t.Errorf("Expected [or]role=admin in query, got %s", result)
	}
}

// ============================================
// 5. NAMED OR GROUPS
// ============================================

func TestOrGroups(t *testing.T) {
	query := cocobase.NewQuery().
		OrGroup("tier").
		Where("isPremium", true).
		Where("isVerified", true).
		Done().
		OrGroup("location").
		Where("country", "US").
		Where("country", "UK").
		Done()

	result := query.Build()

	// Check tier group
	if !strings.Contains(result, "%5Bor%3Atier%5DisPremium=true") {
		t.Errorf("Expected [or:tier]isPremium=true, got %s", result)
	}

	// Check location group
	if !strings.Contains(result, "%5Bor%3Alocation%5Dcountry=US") {
		t.Errorf("Expected [or:location]country=US, got %s", result)
	}
}

// ============================================
// 6. MIXED AND + OR
// ============================================

func TestMixedAndOr(t *testing.T) {
	query := cocobase.NewQuery().
		Where("status", "active").
		GreaterThanOrEqual("age", 18).
		Or().
		Where("isPremium", true).
		Where("isVerified", true).
		Done()

	result := query.Build()
	params := parseQuery(result)

	// Check AND conditions
	if params.Get("status") != "active" {
		t.Errorf("Expected status=active")
	}
	if params.Get("age_gte") != "18" {
		t.Errorf("Expected age_gte=18")
	}

	// Check OR conditions
	if !strings.Contains(result, "%5Bor%5DisPremium=true") {
		t.Errorf("Expected OR condition for isPremium")
	}
}

// ============================================
// 7. PAGINATION
// ============================================

func TestLimit(t *testing.T) {
	query := cocobase.NewQuery().
		Limit(50)

	result := query.Build()

	if !hasParam(result, "limit", "50") {
		t.Errorf("Expected limit=50, got %s", result)
	}
}

func TestOffset(t *testing.T) {
	query := cocobase.NewQuery().
		Offset(100)

	result := query.Build()

	if !hasParam(result, "offset", "100") {
		t.Errorf("Expected offset=100, got %s", result)
	}
}

func TestPage(t *testing.T) {
	query := cocobase.NewQuery().
		Page(3, 20)

	result := query.Build()
	params := parseQuery(result)

	if params.Get("limit") != "20" {
		t.Errorf("Expected limit=20, got %s", params.Get("limit"))
	}
	// Page 3, 20 per page = offset 40
	if params.Get("offset") != "40" {
		t.Errorf("Expected offset=40, got %s", params.Get("offset"))
	}
}

func TestPageWithInvalidNumber(t *testing.T) {
	query := cocobase.NewQuery().
		Page(0, 20)

	result := query.Build()
	params := parseQuery(result)

	// Should default to page 1
	if params.Get("offset") != "0" {
		t.Errorf("Expected offset=0 for invalid page, got %s", params.Get("offset"))
	}
}

// ============================================
// 8. SORTING
// ============================================

func TestOrderBy(t *testing.T) {
	query := cocobase.NewQuery().
		OrderBy("createdAt")

	result := query.Build()
	params := parseQuery(result)

	if params.Get("sort") != "createdAt" {
		t.Errorf("Expected sort=createdAt")
	}
	if params.Get("order") != "asc" {
		t.Errorf("Expected order=asc (default)")
	}
}

func TestOrderByAsc(t *testing.T) {
	query := cocobase.NewQuery().
		OrderByAsc("age")

	result := query.Build()
	params := parseQuery(result)

	if params.Get("sort") != "age" {
		t.Errorf("Expected sort=age")
	}
	if params.Get("order") != "asc" {
		t.Errorf("Expected order=asc")
	}
}

func TestOrderByDesc(t *testing.T) {
	query := cocobase.NewQuery().
		OrderByDesc("createdAt")

	result := query.Build()
	params := parseQuery(result)

	if params.Get("sort") != "createdAt" {
		t.Errorf("Expected sort=createdAt")
	}
	if params.Get("order") != "desc" {
		t.Errorf("Expected order=desc")
	}
}

// ============================================
// 9. HELPER METHODS
// ============================================

func TestActive(t *testing.T) {
	query := cocobase.NewQuery().
		Active()

	result := query.Build()

	if !hasParam(result, "deletedAt_isnull", "true") {
		t.Errorf("Expected deletedAt_isnull=true for Active(), got %s", result)
	}
}

func TestDeleted(t *testing.T) {
	query := cocobase.NewQuery().
		Deleted()

	result := query.Build()

	if !hasParam(result, "deletedAt_isnull", "false") {
		t.Errorf("Expected deletedAt_isnull=false for Deleted(), got %s", result)
	}
}

func TestRecent(t *testing.T) {
	query := cocobase.NewQuery().
		Recent()

	result := query.Build()
	params := parseQuery(result)

	if params.Get("sort") != "created_at" {
		t.Errorf("Expected sort=created_at for Recent()")
	}
	if params.Get("order") != "desc" {
		t.Errorf("Expected order=desc for Recent()")
	}
}

func TestOldest(t *testing.T) {
	query := cocobase.NewQuery().
		Oldest()

	result := query.Build()
	params := parseQuery(result)

	if params.Get("sort") != "created_at" {
		t.Errorf("Expected sort=created_at for Oldest()")
	}
	if params.Get("order") != "asc" {
		t.Errorf("Expected order=asc for Oldest()")
	}
}

// ============================================
// 10. COMPLEX REAL-WORLD QUERIES
// ============================================

func TestComplexEcommerce(t *testing.T) {
	query := cocobase.NewQuery().
		Where("inStock", true).
		Between("price", 10, 100).
		Or().
		Where("onSale", true).
		Where("isNew", true).
		Done().
		OrderBy("price").
		Limit(50)

	result := query.Build()
	params := parseQuery(result)

	// Check all components
	if params.Get("inStock") != "true" {
		t.Errorf("Expected inStock=true")
	}
	if params.Get("price_gte") != "10" || params.Get("price_lte") != "100" {
		t.Errorf("Expected price between 10 and 100")
	}
	if params.Get("limit") != "50" {
		t.Errorf("Expected limit=50")
	}
	if params.Get("sort") != "price" {
		t.Errorf("Expected sort=price")
	}
}

func TestComplexUserSearch(t *testing.T) {
	query := cocobase.NewQuery().
		Active().
		OrGroup("tier").
		Where("isPremium", true).
		Where("isVerified", true).
		Done().
		OrGroup("location").
		Where("country", "US").
		Where("country", "UK").
		Done().
		Recent().
		Page(2, 20)

	result := query.Build()
	params := parseQuery(result)

	// Check active users
	if params.Get("deletedAt_isnull") != "true" {
		t.Errorf("Expected active users filter")
	}

	// Check pagination
	if params.Get("limit") != "20" || params.Get("offset") != "20" {
		t.Errorf("Expected page 2 with 20 per page")
	}

	// Check sorting
	if params.Get("sort") != "created_at" || params.Get("order") != "desc" {
		t.Errorf("Expected recent ordering")
	}
}

func TestComplexTaskQuery(t *testing.T) {
	query := cocobase.NewQuery().
		Where("assignedTo", "user123").
		NotIn("status", "completed", "cancelled").
		Or().
		Equals("priority", "high").
		Where("isOverdue", true).
		Done().
		OrderByDesc("priority").
		Limit(100)

	result := query.Build()
	params := parseQuery(result)

	if params.Get("assignedTo") != "user123" {
		t.Errorf("Expected assignedTo=user123")
	}
	if params.Get("status_notin") != "completed,cancelled" {
		t.Errorf("Expected status_notin")
	}
	if params.Get("limit") != "100" {
		t.Errorf("Expected limit=100")
	}
}

// ============================================
// 11. CHAINING TESTS
// ============================================

func TestMethodChaining(t *testing.T) {
	query := cocobase.NewQuery().
		Where("status", "active").
		GreaterThanOrEqual("age", 18).
		Contains("name", "john").
		In("role", "admin", "user").
		Active().
		Recent().
		Page(1, 20)

	result := query.Build()
	params := parseQuery(result)

	// Verify all chained methods worked
	if params.Get("status") != "active" {
		t.Errorf("Chaining failed: status")
	}
	if params.Get("age_gte") != "18" {
		t.Errorf("Chaining failed: age_gte")
	}
	if params.Get("name_contains") != "john" {
		t.Errorf("Chaining failed: name_contains")
	}
	if params.Get("limit") != "20" {
		t.Errorf("Chaining failed: limit")
	}
}

// ============================================
// 12. EDGE CASES
// ============================================

func TestEmptyQuery(t *testing.T) {
	query := cocobase.NewQuery()
	result := query.Build()

	if result != "" {
		t.Errorf("Expected empty query string, got %s", result)
	}
}

func TestMultipleOrGroups(t *testing.T) {
	query := cocobase.NewQuery().
		OrGroup("group1").
		Where("a", 1).
		Where("b", 2).
		Done().
		OrGroup("group2").
		Where("c", 3).
		Where("d", 4).
		Done().
		OrGroup("group3").
		Where("e", 5).
		Done()

	result := query.Build()

	// Should contain all three groups
	if !strings.Contains(result, "%5Bor%3Agroup1%5D") {
		t.Errorf("Missing group1")
	}
	if !strings.Contains(result, "%5Bor%3Agroup2%5D") {
		t.Errorf("Missing group2")
	}
	if !strings.Contains(result, "%5Bor%3Agroup3%5D") {
		t.Errorf("Missing group3")
	}
}

func TestOrBuilderChaining(t *testing.T) {
	query := cocobase.NewQuery().
		Or().
		Where("a", 1).
		GreaterThan("b", 2).
		Contains("c", "test").
		IsNull("d").
		Done()

	result := query.Build()

	// All OR conditions should be present
	if !strings.Contains(result, "%5Bor%5Da=1") {
		t.Errorf("Missing OR condition: a=1")
	}
	if !strings.Contains(result, "%5Bor%5Db_gt=2") {
		t.Errorf("Missing OR condition: b_gt=2")
	}
	if !strings.Contains(result, "%5Bor%5Dc_contains=test") {
		t.Errorf("Missing OR condition: c_contains=test")
	}
}

// ============================================
// BENCHMARK TESTS
// ============================================

func BenchmarkSimpleQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cocobase.NewQuery().
			Where("status", "active").
			Limit(10).
			Build()
	}
}

func BenchmarkComplexQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cocobase.NewQuery().
			Where("status", "active").
			Between("age", 18, 65).
			Or().
			Where("isPremium", true).
			Where("isVerified", true).
			Done().
			Recent().
			Page(2, 20).
			Build()
	}
}

func BenchmarkOrGroups(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cocobase.NewQuery().
			OrGroup("tier").
			Where("isPremium", true).
			Where("isVerified", true).
			Done().
			OrGroup("location").
			Where("country", "US").
			Where("country", "UK").
			Done().
			Build()
	}
}
