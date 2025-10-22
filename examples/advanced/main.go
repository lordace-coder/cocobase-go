package main

import (
	"context"
	"fmt"

	"github.com/lordace-coder/cocobase-go/cocobase"
)

func main() {
	client := cocobase.NewClient(cocobase.Config{
		APIKey: "your-api-key",
	})

	ctx := context.Background()

	// Run examples
	fmt.Println("🚀 Cocobase Go Client - Intuitive Query Examples\n")

	basicExamples(ctx, client)
	comparisonExamples(ctx, client)
	stringExamples(ctx, client)
	listExamples(ctx, client)
	booleanLogicExamples(ctx, client)
	paginationExamples(ctx, client)
	realWorldExamples(ctx, client)
}

// ============================================
// 1. BASIC QUERIES (Super Simple)
// ============================================

func basicExamples(ctx context.Context, client *cocobase.Client) {
	fmt.Println("=== 1. Basic Queries ===\n")

	// Find active users
	query := cocobase.NewQuery().
		Where("status", "active")

	docs, _ := client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Active users: %d\n", len(docs))

	// Find users aged 18-65
	query = cocobase.NewQuery().
		Between("age", 18, 65)

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Users aged 18-65: %d\n", len(docs))

	// Find active, non-deleted users
	query = cocobase.NewQuery().
		Where("status", "active").
		Active() // Helper method: deletedAt IS NULL

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Active non-deleted users: %d\n\n", len(docs))
}

// ============================================
// 2. COMPARISON OPERATORS (Intuitive Names)
// ============================================

func comparisonExamples(ctx context.Context, client *cocobase.Client) {
	fmt.Println("=== 2. Comparison Operators ===\n")

	// Adults (age >= 18)
	query := cocobase.NewQuery().GreaterThanOrEqual("age", 18)

	fmt.Println("✓ GreaterThanOrEqual(age, 18)", query)

	// Seniors (age > 65)
	query = cocobase.NewQuery().
		GreaterThan("age", 65)

	fmt.Println("✓ GreaterThan(age, 65)")

	// Not banned users
	query = cocobase.NewQuery().
		NotEquals("status", "banned")

	fmt.Println("✓ NotEquals(status, banned)")

	// Age range: 18-65
	query = cocobase.NewQuery().
		Between("age", 18, 65)

	fmt.Println("✓ Between(age, 18, 65)\n")
}

// ============================================
// 3. STRING OPERATIONS (Easy to Remember)
// ============================================

func stringExamples(ctx context.Context, client *cocobase.Client) {
	fmt.Println("=== 3. String Operations ===\n")

	// Search for "john" in name
	query := cocobase.NewQuery().
		Contains("name", "john")

	docs, _ := client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Contains(name, john): %d users\n", len(docs))

	// Find Gmail users
	query = cocobase.NewQuery().
		EndsWith("email", "gmail.com")

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ EndsWith(email, gmail.com): %d users\n", len(docs))

	// Names starting with "John"
	query = cocobase.NewQuery().
		StartsWith("name", "John")

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ StartsWith(name, John): %d users\n", len(docs))

	// Search "admin" in name OR email OR username
	query = cocobase.NewQuery().
		Search("admin", "name", "email", "username")

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Search(admin) in 3 fields: %d users\n\n", len(docs))
}

// ============================================
// 4. LIST OPERATIONS (Simple & Clean)
// ============================================

func listExamples(ctx context.Context, client *cocobase.Client) {
	fmt.Println("=== 4. List Operations ===\n")

	// Find admins, moderators, or support staff
	query := cocobase.NewQuery().
		In("role", "admin", "moderator", "support")

	docs, _ := client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ In(role, admin/moderator/support): %d users\n", len(docs))

	// Exclude banned and deleted users
	query = cocobase.NewQuery().
		NotIn("status", "banned", "deleted", "suspended")

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ NotIn(status, banned/deleted/suspended): %d users\n", len(docs))

	// Users with profile pictures
	query = cocobase.NewQuery().
		IsNotNull("profilePicture")

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ IsNotNull(profilePicture): %d users\n\n", len(docs))
}

// ============================================
// 5. BOOLEAN LOGIC (Super Intuitive)
// ============================================

func booleanLogicExamples(ctx context.Context, client *cocobase.Client) {
	fmt.Println("=== 5. Boolean Logic (OR Conditions) ===\n")

	// Simple OR: Premium OR Verified users
	query := cocobase.NewQuery().
		Or().
		Where("isPremium", true).
		Where("isVerified", true).
		Done()

	docs, _ := client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Premium OR Verified: %d users\n", len(docs))

	// Mixed AND + OR: Active AND (Premium OR Verified)
	query = cocobase.NewQuery().
		Where("status", "active").
		Or().
		Where("isPremium", true).
		Where("isVerified", true).
		Done()

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Active AND (Premium OR Verified): %d users\n", len(docs))

	// Multiple OR Groups: (Premium OR Verified) AND (US OR UK)
	query = cocobase.NewQuery().
		OrGroup("tier").
		Where("isPremium", true).
		Where("isVerified", true).
		Done().
		OrGroup("location").
		Where("country", "US").
		Where("country", "UK").
		Done()

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Multiple OR groups: %d users\n", len(docs))

	// Complex: Active AND Age>=18 AND (Admin OR Moderator)
	query = cocobase.NewQuery().
		Where("status", "active").
		GreaterThanOrEqual("age", 18).
		Or().
		Equals("role", "admin").
		Equals("role", "moderator").
		Done()

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Complex query: %d users\n\n", len(docs))
}

// ============================================
// 6. PAGINATION & SORTING (Dead Simple)
// ============================================

func paginationExamples(ctx context.Context, client *cocobase.Client) {
	fmt.Println("=== 6. Pagination & Sorting ===\n")

	// Most recent 10 users
	query := cocobase.NewQuery().
		Recent().
		Limit(10)

	docs, _ := client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Recent().Limit(10): %d users\n", len(docs))

	// Oldest users first
	query = cocobase.NewQuery().
		Oldest().
		Limit(20)

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Oldest().Limit(20): %d users\n", len(docs))

	// Custom sort: by age descending
	query = cocobase.NewQuery().
		OrderByDesc("age").
		Limit(50)

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ OrderByDesc(age).Limit(50): %d users\n", len(docs))

	// Page 2 (20 per page)
	query = cocobase.NewQuery().
		Page(2, 20).
		Recent()

	docs, _ = client.ListDocuments(ctx, "users", query)
	fmt.Printf("✓ Page(2, 20): %d users\n\n", len(docs))
}

// ============================================
// 7. REAL-WORLD EXAMPLES
// ============================================

func realWorldExamples(ctx context.Context, client *cocobase.Client) {
	fmt.Println("=== 7. Real-World Use Cases ===\n")

	// 1. E-commerce: Available products on sale
	query := cocobase.NewQuery().
		Where("inStock", true).
		Where("onSale", true).
		Between("price", 10, 100).
		OrderBy("price").
		Limit(50)

	fmt.Println("✓ E-commerce: In-stock sale items $10-$100 ", query)

	// 2. Social Media: Popular recent posts
	query = cocobase.NewQuery().
		Recent().
		Or().
		GreaterThan("likes", 100).
		GreaterThan("comments", 50).
		Done().
		IsNull("deletedAt").
		Limit(20)

	fmt.Println("✓ Social: Popular recent posts")

	// 3. User Management: Find risky accounts
	query = cocobase.NewQuery().
		Or().
		GreaterThanOrEqual("failedLogins", 5).
		Where("suspiciousActivity", true).
		Done().
		NotEquals("status", "banned").
		Recent()

	fmt.Println("✓ Security: Risky user accounts")

	// 4. Task Management: My urgent tasks
	query = cocobase.NewQuery().
		Where("assignedTo", "user123").
		NotIn("status", "completed", "cancelled").
		Or().
		Equals("priority", "high").
		Where("isOverdue", true).
		Done().
		OrderByDesc("priority")

	fmt.Println("✓ Tasks: My urgent incomplete tasks")

	// 5. Content Moderation: Flagged items
	query = cocobase.NewQuery().
		Where("status", "pending_review").
		Or().
		GreaterThanOrEqual("reports", 3).
		Contains("content", "spam").
		Done().
		Recent().
		Limit(100)

	fmt.Println("✓ Moderation: Flagged content")

	// 6. Analytics: Active premium users this month
	query = cocobase.NewQuery().
		Where("isPremium", true).
		Active().
		GreaterThanOrEqual("lastLoginAt", "2025-01-01").
		OrderByDesc("lastLoginAt")

	fmt.Println("✓ Analytics: Active premium users")

	// 7. Support: Open tickets from VIP customers
	query = cocobase.NewQuery().
		In("status", "open", "pending", "escalated").
		OrGroup("priority").
		Where("isVIP", true).
		Equals("priority", "high").
		Done().
		Recent()

	fmt.Println("✓ Support: VIP or high-priority tickets")

	// 8. Inventory: Low stock alerts
	query = cocobase.NewQuery().
		Or().
		LessThanOrEqual("stock", 10).
		Equals("stock", 0).
		Done().
		Where("discontinued", false).
		OrderBy("stock")

	fmt.Println("✓ Inventory: Low stock items\n")
}

// ============================================
// COMPARISON: Old Way vs New Way
// ============================================

func comparisonOldVsNew() {
	fmt.Println("=== 📊 Old Way vs New Way ===\n")

	fmt.Println("❌ OLD (Hard to remember):")
	fmt.Println(`query := cocobase.NewQuery().
    Filter("age", "gte", 18).
    Filter("email", "endswith", "gmail.com").
    Or("isPremium", "eq", true).
    Or("isVerified", "eq", true)`)

	fmt.Println("\n✅ NEW (Intuitive):")
	fmt.Println(`query := cocobase.NewQuery().
    GreaterThanOrEqual("age", 18).
    EndsWith("email", "gmail.com").
    Or().
        Where("isPremium", true).
        Where("isVerified", true).
    Done()`)

	fmt.Println("\n🎯 Much easier to read and remember!")
}
