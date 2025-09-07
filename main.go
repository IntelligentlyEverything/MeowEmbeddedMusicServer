package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	TAG = "MeowEmbeddedMusicServer"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("[Warning] %s Loading .env file failed: %v\nUse the default configuration instead.\n", TAG, err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Printf("[Warning] %s PORT environment variable not set\nUse the default port 2233 instead.\n", TAG)
		port = "2233"
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/stream_pcm", apiHandler)
	fmt.Printf("[Info] %s Started.\n喵波音律-音乐家园QQ交流群:865754861\n", TAG)
	fmt.Printf("[Info] Starting music server at port %s\n", port)

	// Create a channel to listen for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a server instance
	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
			sigChan <- syscall.SIGINT // Send a signal to shut down the server
		}
	}()

	// Create a channel to listen for standard input
	exitChan := make(chan struct{})

	go func() {
		for {
			var input string
			fmt.Scanln(&input)
			if input == "exit" {
				exitChan <- struct{}{}
				return
			}
		}
	}()

	// Monitor signals or exit signals from standard inputs
	select {
	case <-sigChan:
		fmt.Printf("[Info] Server is shutting down.\nGoodbye!\n")
	case <-exitChan:
		fmt.Printf("[Info] Server is shutting down.\nGoodbye!\n")
	}

	// Shut down the server
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Println(err)
	}
}
