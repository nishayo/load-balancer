package algorithms

import (
	"sync/atomic"
)

// Global variables
var counter uint32   // Round-robin counter
var Servers []string // Exported servers slice

// Initialize sets up the servers for the round robin algorithm
func Initialize(serverList []string) {
	Servers = serverList
}

// GetNextServer selects the next backend server using round-robin
func GetNextServer() string {
	// Check if Servers slice is empty to avoid division by zero
	if len(Servers) == 0 {
		return ""
	}
	
	index := atomic.AddUint32(&counter, 1) % uint32(len(Servers))
	return Servers[index]
}