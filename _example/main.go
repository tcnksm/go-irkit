package main

import (
	"context"
	"log"
	"os"

	"github.com/tcnksm/go-irkit/v1"
)

func main() {
	c := irkit.DefaultInternetClient()
	id, key, err := c.GetKeys(context.Background(), os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	log.Printf("[INFO] id: %s", id)
	log.Printf("[INFO] key: %s", key)
}
