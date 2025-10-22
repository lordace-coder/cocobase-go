package cocobase

import (
	"fmt"
	"net/url"
	"strings"
)

// QueryBuilder provides a fluent, intuitive interface for building queries
type QueryBuilder struct {
	filters   map[string]string
	orFilters map[string][]string
	limit     int
	offset    int
	sort      string
	order     string
}

// NewQuery creates a new QueryBuilder
func NewQuery() *QueryBuilder {
	return &QueryBuilder{
		filters:   make(map[string]string),
		orFilters: make(map[string][]string),
	}
}

// ============================================
// COMPARISON OPERATORS (Intuitive Names)
// ============================================

// Where adds an equality filter (field = value)
func (qb *QueryBuilder) Where(field string, value interface{}) *QueryBuilder {
	qb.filters[field] = fmt.Sprintf("%v", value)
	return qb
}

// Equals is an alias for Where
func (qb *QueryBuilder) Equals(field string, value interface{}) *QueryBuilder {
	return qb.Where(field, value)
}

// NotEquals adds a not-equals filter (field != value)
func (qb *QueryBuilder) NotEquals(field string, value interface{}) *QueryBuilder {
	key := fmt.Sprintf("%s_ne", field)
	qb.filters[key] = fmt.Sprintf("%v", value)
	return qb
}

// GreaterThan adds a greater-than filter (field > value)
func (qb *QueryBuilder) GreaterThan(field string, value interface{}) *QueryBuilder {
	key := fmt.Sprintf("%s_gt", field)
	qb.filters[key] = fmt.Sprintf("%v", value)
	return qb
}

// GreaterThanOrEqual adds a gte filter (field >= value)
func (qb *QueryBuilder) GreaterThanOrEqual(field string, value interface{}) *QueryBuilder {
	key := fmt.Sprintf("%s_gte", field)
	qb.filters[key] = fmt.Sprintf("%v", value)
	return qb
}

// LessThan adds a less-than filter (field < value)
func (qb *QueryBuilder) LessThan(field string, value interface{}) *QueryBuilder {
	key := fmt.Sprintf("%s_lt", field)
	qb.filters[key] = fmt.Sprintf("%v", value)
	return qb
}

// LessThanOrEqual adds a lte filter (field <= value)
func (qb *QueryBuilder) LessThanOrEqual(field string, value interface{}) *QueryBuilder {
	key := fmt.Sprintf("%s_lte", field)
	qb.filters[key] = fmt.Sprintf("%v", value)
	return qb
}

// Between adds a range filter (field >= min AND field <= max)
func (qb *QueryBuilder) Between(field string, min, max interface{}) *QueryBuilder {
	qb.GreaterThanOrEqual(field, min)
	qb.LessThanOrEqual(field, max)
	return qb
}

// ============================================
// STRING OPERATORS (Easy to Remember)
// ============================================

// Contains adds a substring search filter (case-insensitive)
func (qb *QueryBuilder) Contains(field, substring string) *QueryBuilder {
	key := fmt.Sprintf("%s_contains", field)
	qb.filters[key] = substring
	return qb
}

// StartsWith adds a prefix filter
func (qb *QueryBuilder) StartsWith(field, prefix string) *QueryBuilder {
	key := fmt.Sprintf("%s_startswith", field)
	qb.filters[key] = prefix
	return qb
}

// EndsWith adds a suffix filter
func (qb *QueryBuilder) EndsWith(field, suffix string) *QueryBuilder {
	key := fmt.Sprintf("%s_endswith", field)
	qb.filters[key] = suffix
	return qb
}

// Search searches across multiple fields (multi-field OR)
func (qb *QueryBuilder) Search(searchTerm string, fields ...string) *QueryBuilder {
	key := strings.Join(fields, "__or__") + "_contains"
	qb.filters[key] = searchTerm
	return qb
}

// ============================================
// LIST OPERATORS
// ============================================

// In adds an "in list" filter
func (qb *QueryBuilder) In(field string, values ...interface{}) *QueryBuilder {
	key := fmt.Sprintf("%s_in", field)
	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = fmt.Sprintf("%v", v)
	}
	qb.filters[key] = strings.Join(strValues, ",")
	return qb
}

// NotIn adds a "not in list" filter
func (qb *QueryBuilder) NotIn(field string, values ...interface{}) *QueryBuilder {
	key := fmt.Sprintf("%s_notin", field)
	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = fmt.Sprintf("%v", v)
	}
	qb.filters[key] = strings.Join(strValues, ",")
	return qb
}

// ============================================
// NULL CHECKS
// ============================================

// IsNull adds a null check filter
func (qb *QueryBuilder) IsNull(field string) *QueryBuilder {
	key := fmt.Sprintf("%s_isnull", field)
	qb.filters[key] = "true"
	return qb
}

// IsNotNull adds a not-null check filter
func (qb *QueryBuilder) IsNotNull(field string) *QueryBuilder {
	key := fmt.Sprintf("%s_isnull", field)
	qb.filters[key] = "false"
	return qb
}

// ============================================
// BOOLEAN LOGIC (Simple OR)
// ============================================

// Or starts a simple OR condition
func (qb *QueryBuilder) Or() *OrBuilder {
	return &OrBuilder{
		qb:        qb,
		groupName: "",
	}
}

// OrGroup starts a named OR group
func (qb *QueryBuilder) OrGroup(groupName string) *OrBuilder {
	return &OrBuilder{
		qb:        qb,
		groupName: groupName,
	}
}

// OrBuilder provides fluent OR condition building
type OrBuilder struct {
	qb        *QueryBuilder
	groupName string
}

// Where adds an equality condition to the OR
func (ob *OrBuilder) Where(field string, value interface{}) *OrBuilder {
	return ob.addCondition(field, "", value)
}

// Equals is an alias for Where
func (ob *OrBuilder) Equals(field string, value interface{}) *OrBuilder {
	return ob.Where(field, value)
}

// NotEquals adds a not-equals OR condition
func (ob *OrBuilder) NotEquals(field string, value interface{}) *OrBuilder {
	return ob.addCondition(field, "ne", value)
}

// GreaterThan adds a greater-than OR condition
func (ob *OrBuilder) GreaterThan(field string, value interface{}) *OrBuilder {
	return ob.addCondition(field, "gt", value)
}

// GreaterThanOrEqual adds a gte OR condition
func (ob *OrBuilder) GreaterThanOrEqual(field string, value interface{}) *OrBuilder {
	return ob.addCondition(field, "gte", value)
}

// LessThan adds a less-than OR condition
func (ob *OrBuilder) LessThan(field string, value interface{}) *OrBuilder {
	return ob.addCondition(field, "lt", value)
}

// LessThanOrEqual adds a lte OR condition
func (ob *OrBuilder) LessThanOrEqual(field string, value interface{}) *OrBuilder {
	return ob.addCondition(field, "lte", value)
}

// Contains adds a contains OR condition
func (ob *OrBuilder) Contains(field, substring string) *OrBuilder {
	return ob.addCondition(field, "contains", substring)
}

// StartsWith adds a startswith OR condition
func (ob *OrBuilder) StartsWith(field, prefix string) *OrBuilder {
	return ob.addCondition(field, "startswith", prefix)
}

// EndsWith adds an endswith OR condition
func (ob *OrBuilder) EndsWith(field, suffix string) *OrBuilder {
	return ob.addCondition(field, "endswith", suffix)
}

// IsNull adds a null check OR condition
func (ob *OrBuilder) IsNull(field string) *OrBuilder {
	return ob.addCondition(field, "isnull", true)
}

// IsNotNull adds a not-null check OR condition
func (ob *OrBuilder) IsNotNull(field string) *OrBuilder {
	return ob.addCondition(field, "isnull", false)
}

// Done finishes the OR builder and returns the main QueryBuilder
func (ob *OrBuilder) Done() *QueryBuilder {
	return ob.qb
}

func (ob *OrBuilder) addCondition(field, operator string, value interface{}) *OrBuilder {
	key := field
	if operator != "" {
		key = fmt.Sprintf("%s_%s", field, operator)
	}

	var prefix string
	if ob.groupName == "" {
		prefix = "[or]"
	} else {
		prefix = fmt.Sprintf("[or:%s]", ob.groupName)
	}

	filterStr := fmt.Sprintf("%s%s=%v", prefix, key, value)
	ob.qb.orFilters[ob.groupName] = append(ob.qb.orFilters[ob.groupName], filterStr)
	return ob
}

// ============================================
// PAGINATION
// ============================================

// Limit sets the maximum number of results
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset sets the number of results to skip
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Page sets pagination using page number (starts at 1)
func (qb *QueryBuilder) Page(page, perPage int) *QueryBuilder {
	if page < 1 {
		page = 1
	}
	qb.limit = perPage
	qb.offset = (page - 1) * perPage
	return qb
}

// ============================================
// SORTING
// ============================================

// OrderBy sets the field to sort by (ascending by default)
func (qb *QueryBuilder) OrderBy(field string) *QueryBuilder {
	qb.sort = field
	qb.order = "asc"
	return qb
}

// OrderByAsc sets ascending sort order
func (qb *QueryBuilder) OrderByAsc(field string) *QueryBuilder {
	qb.sort = field
	qb.order = "asc"
	return qb
}

// OrderByDesc sets descending sort order
func (qb *QueryBuilder) OrderByDesc(field string) *QueryBuilder {
	qb.sort = field
	qb.order = "desc"
	return qb
}

// Asc sets ascending order (use after OrderBy)
func (qb *QueryBuilder) Asc() *QueryBuilder {
	qb.order = "asc"
	return qb
}

// Desc sets descending order (use after OrderBy)
func (qb *QueryBuilder) Desc() *QueryBuilder {
	qb.order = "desc"
	return qb
}

// ============================================
// BUILD QUERY STRING
// ============================================

// Build constructs the final query string
func (qb *QueryBuilder) Build() string {
	params := url.Values{}

	// Add simple AND filters
	for key, value := range qb.filters {
		params.Add(key, value)
	}

	// Add OR filters (all groups)
	for _, filters := range qb.orFilters {
		for _, filter := range filters {
			parts := strings.SplitN(filter, "=", 2)
			if len(parts) == 2 {
				params.Add(parts[0], parts[1])
			}
		}
	}

	// Add pagination
	if qb.limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", qb.limit))
	}
	if qb.offset > 0 {
		params.Add("offset", fmt.Sprintf("%d", qb.offset))
	}

	// Add sorting
	if qb.sort != "" {
		params.Add("sort", qb.sort)
		if qb.order != "" {
			params.Add("order", qb.order)
		}
	}

	return params.Encode()
}

// ============================================
// HELPER METHODS FOR COMMON PATTERNS
// ============================================

// Active adds a filter for active/non-deleted records
func (qb *QueryBuilder) Active() *QueryBuilder {
	return qb.IsNull("deletedAt")
}

// Deleted adds a filter for deleted records
func (qb *QueryBuilder) Deleted() *QueryBuilder {
	return qb.IsNotNull("deletedAt")
}

// Recent orders by created_at descending
func (qb *QueryBuilder) Recent() *QueryBuilder {
	return qb.OrderByDesc("created_at")
}

// Oldest orders by created_at ascending
func (qb *QueryBuilder) Oldest() *QueryBuilder {
	return qb.OrderByAsc("created_at")
}
