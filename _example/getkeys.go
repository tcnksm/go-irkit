package main

import (
	"context"
	"log"
	"os"

	"github.com/tcnksm/go-irkit/v1"
)

func main() {
	c := irkit.DefaultInternetClient()
	deviceid, clientkey, err := c.GetKeys(context.Background(), os.Getenv("CLIENT_TOKEN"))
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	log.Printf("[INFO] deviceid: %s", deviceid)
	log.Printf("[INFO] clientkey: %s", clientkey)
}
