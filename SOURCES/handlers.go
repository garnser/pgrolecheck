// handlers.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func CheckRoleHandler(w http.ResponseWriter, cfg Configuration) {
	// Database interaction moved to db.go, call it here
	status, err := checkDatabaseRole(cfg)
	if err != nil {
		respondWithError(w, "notok")
		log.Println("Error:", err)
		return
	}
	respondWithStatus(w, status)
}

func respondWithError(w http.ResponseWriter, status string) {
   respondWithStatus(w, status)
}

func respondWithStatus(w http.ResponseWriter, status string) {
   response := map[string]string{"status": status}
   jsonResponse, err := json.Marshal(response)
   if err != nil {
       log.Println("Error creating JSON response:", err)
       http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
       return
   }
   
   w.Header().Set("Content-Type", "application/json")
   w.WriteHeader(http.StatusOK) // Consider adjusting the status code based on the error
   w.Write(jsonResponse)
}


