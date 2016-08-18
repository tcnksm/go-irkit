package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/tcnksm/go-irkit/v1"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("[ERROR] missing arguments")
	}

	filePath := os.Args[1]
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	var msg irkit.Message
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&msg); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	c := irkit.DefaultInternetClient()
	err = c.SendMessages(context.Background(), os.Getenv("CLIENT_KEY"), os.Getenv("DEVICE_ID"), &msg)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
