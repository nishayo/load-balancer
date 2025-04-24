package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func SetupConfig() {
	reader := os.Stdin
	fmt.Println("\n\033[1m⚙️  Load Balancer Configuration\033[0m\n") // Bold title

	// Algorithm selection
	fmt.Println("\033[1mSelect Load Balancing Algorithm:\033[0m")
	fmt.Println("1) Round-Robin")
	fmt.Println("2) Least Connections")
	fmt.Println("3) IP Hashing")
	fmt.Println("4) Latency-Based")
	fmt.Print("Enter choice (1-4): ")
	var algoChoice string
	fmt.Fscanln(reader, &algoChoice)

	algorithms := map[string]string{
		"1": "round_robin",
		"2": "least_connections",
		"3": "ip_hashing",
		"4": "latency_based",
	}
	algorithm, exists := algorithms[algoChoice]
	if !exists {
		fmt.Println("Invalid choice, defaulting to Round-Robin")
		algorithm = "round_robin"
	}

	// Backend server selection
	fmt.Println("\n\033[1mEnter backend servers (comma-separated, or press Enter for default [http://localhost:5001, http://localhost:5002, http://localhost:5003]):\033[0m")
	var serversInput string
	fmt.Fscanln(reader, &serversInput)

	defaultServers := []string{"http://localhost:5001", "http://localhost:5002", "http://localhost:5003"}
	servers := defaultServers
	if serversInput != "" {
		servers = strings.Split(serversInput, ",")
	}

	// Start default servers if user didn't provide custom ones
	if serversInput == "" {
		fmt.Println("Using default servers.")
		go MyServers()
	}

	// Sticky sessions
	fmt.Print("\n\033[1mEnable Sticky Sessions? (y/n): \033[0m")
	var stickyInput string
	fmt.Fscanln(reader, &stickyInput)
	stickySession := strings.ToLower(stickyInput) == "y"

	// Save config to TOML
	config := Config{Algorithm: algorithm, StickySession: stickySession, Servers: servers}
	configData, err := toml.Marshal(config)
	if err != nil {
		fmt.Println("Error encoding TOML:", err)
		return
	}

	err = os.WriteFile("config.toml", configData, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}

	// Display final configurations
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println("\n✅ Configurations saved to config.toml file.")
	fmt.Println(string(configJSON))
}
