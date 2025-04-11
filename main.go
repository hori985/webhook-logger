package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func logJSON(data map[string]interface{}) {
	entry, err := json.Marshal(data)
	if err != nil {
		log.Println(`{"level":"error","message":"Failed to marshal log entry"}`)
		return
	}
	log.Println(string(entry))
}

func handler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logJSON(map[string]any{
			"level":   "error",
			"message": "Failed to read request body",
			"error":   err.Error(),
		})
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	logJSON(map[string]any{
		"level":      "info",
		"method":     r.Method,
		"path":       r.URL.Path,
		"user-agent": r.UserAgent(),
		"body":       string(bodyBytes),
	})

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("POST /", handler)

	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)

	// This will block and run the server
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
