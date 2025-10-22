package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/cocobase-go/cocobase"
	"github.com/yourusername/cocobase-go/storage"
)

func main() {
	// Create client with persistent storage
	store := storage.NewMemoryStorage()
	
	client := cocobase.NewClient(cocobase.Config{
		APIKey:  "your-api-key",
		Storage: store,
	})

	ctx := context.Background()

	// Register a new user
	err := client.Register(ctx, "user@example.com", "password123", map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User registered successfully")

	// Check if authenticated
	if client.IsAuthenticated() {
		fmt.Println("User is authenticated")
	}

	// Get current user
	user, err := client.GetCurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current user: %s\n", user.Email)

	// Update user
	newEmail := "newemail@example.com"
	user, err = client.UpdateUser(ctx, map[string]interface{}{
		"phone": "+1234567890",
	}, &newEmail, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated user: %s\n", user.Email)

	// Logout
	if err := client.Logout(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged out")

	// Login
	err = client.Login(ctx, "newemail@example.com", "password123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in successfully")
}
