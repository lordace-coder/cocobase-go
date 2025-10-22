package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/cocobase-go/cocobase"
)

func main() {
	client := cocobase.NewClient(cocobase.Config{
		APIKey: "your-api-key",
	})

	ctx := context.Background()

	// Create a document
	doc, err := client.CreateDocument(ctx, "users", map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created document: %s\n", doc.ID)

	// Get a document
	doc, err = client.GetDocument(ctx, "users", doc.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved: %+v\n", doc.Data)

	// List documents with simple query
	query := cocobase.NewQuery().
		Where("age", 30).
		Limit(10)

	docs, err := client.ListDocuments(ctx, "users", query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d users\n", len(docs))
}
