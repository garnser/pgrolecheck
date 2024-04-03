// middleware.go

package main

import (
    "log"
    "time"
    "net"
    "net/http"
    "strings"
)

// IsIPAllowed checks if an IP address is in the whitelist. If the whitelist is empty, all IPs are allowed.
func IsIPAllowed(ipAddr string, whitelist []string) bool {
    // If the whitelist is empty, allow all IPs
    if len(whitelist) == 0 {
        return true
    }

    // Parse the source IP address
    srcIP := net.ParseIP(ipAddr)
    if srcIP == nil {
        // If the source IP address is invalid, it's not allowed
        return false
    }

    // Check each entry in the whitelist
    for _, allowed := range whitelist {
        if strings.Contains(allowed, "/") {
            // The entry is in CIDR notation
            _, cidrNet, err := net.ParseCIDR(allowed)
            if err != nil {
                // If the CIDR notation is invalid, skip this entry
                continue
            }
            if cidrNet.Contains(srcIP) {
                return true
            }
        } else {
            // The entry is a plain IP address
            if ipAddr == allowed {
                return true
            }
        }
    }

    // IP is not in the whitelist
    log.Printf("IP Whitelist: %s isn't permitted to connect to the service", srcIP)
    return false
}

// IPWhitelistMiddleware wraps the HTTP handler to allow only requests from whitelisted IPs.
func IPWhitelistMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        clientIP := r.RemoteAddr
        // If using standard http library, RemoteAddr includes the port, so we might need to strip it
        clientIP = strings.Split(clientIP, ":")[0]

        if !IsIPAllowed(clientIP, cfg.IPWhitelist) {
            // Reject the request
            http.Error(w, "Access denied", http.StatusForbidden)
            return
        }

        // IP is allowed, proceed with the next handler
        next(w, r)
    }
}

func TokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // If the AuthToken is not set, bypass the token check
        if cfg.AuthToken != "" {
            // Retrieve the token from the request headers
            token := r.Header.Get("Authorization")

            // Prepare the expected token prefix
            expectedToken := "Bearer " + cfg.AuthToken

            // Check if the token matches the one in the configuration
            if token != expectedToken {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
        }

        // Token is valid or not required, proceed with the next handler
        next(w, r)
    }
}

// LoggingMiddleware logs details about each request including the method, URI, and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("Access log: %s %s %s %v", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
    })
}
