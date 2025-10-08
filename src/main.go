package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/session"
)

// initialize global variables
var apiID int
var apiHash string
var sessionName string
var target int

type Secrets struct {
	ApiID       int    `json:"api_id"`
	ApiHash     string `json:"api_hash"`
	SessionName string `json:"session_name"`
	Target      int    `json:"target"`
	PhoneNumber string `json:"phone_number"`
}
var secrets Secrets

func initClient() {
	// Read from secrets.json in project root
	data, err := os.ReadFile("secrets.json")
	if err != nil {
		log.Fatalf("Error reading secrets.json: %v", err)
		return
	}
	
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		log.Fatalf("Error parsing secrets.json: %v", err)
		return
	}
	
	apiID = secrets.ApiID
	apiHash = secrets.ApiHash
	sessionName = secrets.SessionName
	target = secrets.Target
}

func main() {
	initClient()
	
	// Create session storage using FileStorage 
	sessionStorage := &session.FileStorage{
		Path: filepath.Join(".", sessionName+".json"),
	}
	
	client := telegram.NewClient(apiID, apiHash, telegram.Options{
		SessionStorage: sessionStorage,
	})
	ctx := context.Background()
	
	err := client.Run(ctx, func(ctx context.Context) error {
		// Check if we're already authenticated
		authStatus, err := client.Auth().Status(ctx)
		if err != nil {
			return fmt.Errorf("failed to get auth status: %v", err)
		}
		
		// If not authenticated, perform authentication
		if !authStatus.Authorized {
			fmt.Println("Not authenticated. Starting authentication...")
			
			if err := client.Auth().IfNecessary(ctx, TerminalAuth(secrets.PhoneNumber)); err != nil {
				return fmt.Errorf("auth failed: %v", err)
			}
			fmt.Println("Authentication successful!")
		} else {
			fmt.Println("Already authenticated!")
		}
		
		return nil 
	})
	
	if err != nil {
		log.Fatalf("Error running client: %v", err)
	}
}