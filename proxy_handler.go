package main
import(
	"load_balancer/algorithms"
	"io"
	"net/http"
	"net/url"
	"strings"
	"log"
)	

// GetClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (for clients behind proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// Otherwise use RemoteAddr
	ip := r.RemoteAddr
	// Strip port if present
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

// Decides target url finding function based on algorithm in config.
func GetTargetServer(r *http.Request) *url.URL {
	var targetURL *url.URL
	var err error

	switch config.Algorithm {
	case "round_robin":
		serverURL := algorithms.GetNextServer()
		if serverURL == "" {
			return nil // No servers available
		}
		targetURL, err = url.Parse(serverURL)
	
	default:
		return nil
	}

	if err != nil {
		log.Printf("Error occurred: %v", err)
		return nil
	}
	
	if targetURL == nil { 
		return nil
	}
	
	return targetURL
}

// ProxyHandler handles reverse proxying to backend servers
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the target server to send the request to
	targetURL := GetTargetServer(r)
	if targetURL == nil {
		http.Error(w, "Could not determine target server", http.StatusInternalServerError)
		return
	}

	// Create a new request to the backend
	proxyReq, err := http.NewRequest(r.Method, targetURL.String()+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	proxyReq.Header = r.Header

	// Send the request to the backend server
	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy headers and body to the original response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}