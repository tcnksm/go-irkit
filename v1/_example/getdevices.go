package main

import (
	"context"
	"log"
	"os"

	"github.com/tcnksm/go-irkit/v1"
)

func main() {
	c := irkit.DefaultInternetClient()
	devicekey, deviceid, err := c.GetDevices(context.Background(), os.Getenv("CLIENT_KEY"))
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	log.Printf("[INFO] devicekey: %s", devicekey)
	log.Printf("[INFO] deviceid: %s", deviceid)
}
