package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/pelletier/go-toml/v2"
)

// Config struct
type Config struct {
	Algorithm     string   `toml:"algorithm"`
	StickySession bool     `toml:"sticky_sessions"`
	Servers       []string `toml:"servers"`
}

// Global variables
var servers []string // Backend servers from config
var counter uint32   // Round-robin counter

func main() {
	// Run cli.go before starting to get configs
	fmt.Println("Running configuration setup...\n")
	cmd := exec.Command("go", "run", "cli.go")

	// Show program output in the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // Ensure input works in CLI

	if err := cmd.Run(); err != nil { // Use Run() to wait for the command to finish
		log.Fatal("Error running cli.go:", err)
	}

	// Load configuration from TOML
	config, err := loadConfig("config.toml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Set global variables
	servers = config.Servers

	// Validate servers
	if len(servers) == 0 {
		log.Fatal("No backend servers defined. Check config.toml.")
	}

	// Start the load balancer
	http.HandleFunc("/", ProxyHandler) // Use ProxyHandler from round_robin.go
	log.Println("Load balancer started on :8080 with algorithm:", config.Algorithm, "Sticky Sessions:", config.StickySession)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Load config from a TOML file
func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
