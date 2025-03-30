package main

import (
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
)

// GetNextServer selects the next backend server using round-robin
func GetNextServer() string {
	index := atomic.AddUint32(&counter, 1) % uint32(len(servers))
	return servers[index]
}

// ProxyHandler handles reverse proxying to backend servers
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL, _ := url.Parse(GetNextServer())

	// Create a new request to the backend
	proxyReq, err := http.NewRequest(r.Method, targetURL.String()+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	proxyReq.Header = r.Header

	// Send the request to the backend server
	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers and status code
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}
