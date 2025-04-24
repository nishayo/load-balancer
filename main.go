package main

import (
	"load_balancer/algorithms"
	"log"
	"net/http"
	"os"

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
var config *Config   // Make config globally accessible

func main() {
	// Run the CLI setup to generate config.toml
	SetupConfig()

	// Load configuration from TOML
	var err error
	config, err = loadConfig("config.toml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Set global variables
	servers = config.Servers

	// Validate servers
	if len(servers) == 0 {
		log.Fatal("No backend servers defined. Check config.toml.")
	}

	// Initialize the algorithm with servers from config
	algorithms.Initialize(servers)

	// Start the load balancer
	http.HandleFunc("/", ProxyHandler) // Use ProxyHandler from proxy_handler.go
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
