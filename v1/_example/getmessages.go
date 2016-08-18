package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tcnksm/go-irkit/v1"
)

var timeout = 5 * time.Second

func main() {

	c := irkit.DefaultInternetClient()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultCh := make(chan *irkit.SignalInfo, 1)
	go func() {
		signalInfo, err := c.GetMessages(ctx, os.Getenv("CLIENT_KEY"), true)
		if err != nil {
			fmt.Println("")
			log.Fatalf("[ERROR] %s", err)
		}

		resultCh <- signalInfo
	}()

	fmt.Printf("waiting for signal")
	for {
		ticker := time.NewTicker(1 * time.Second)
		select {
		case <-ticker.C:
			fmt.Printf(".")
		case res := <-resultCh:
			fmt.Println("")
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(res); err != nil {
				log.Fatalf("[ERROR] %s", err)
			}
			return
		}
	}
}
