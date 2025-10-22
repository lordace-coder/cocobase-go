package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/cocobase-go/cocobase"
)

func main() {
	client := cocobase.NewClient(cocobase.Config{
		APIKey: "your-api-key",
	})

	ctx := context.Background()

	// Watch collection for changes
	conn, err := client.WatchCollection(ctx, "users", func(event cocobase.Event) {
		fmt.Printf("Event: %s\n", event.Event)
		fmt.Printf("Document ID: %s\n", event.Data.ID)
		fmt.Printf("Data: %+v\n", event.Data.Data)
		fmt.Println("---")
	}, "users-watcher")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Watching for changes... Press Ctrl+C to exit")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nClosing connection...")
}
