package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"io/ioutil"
)

type Game struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Stores      map[string]string `json:"stores"`
	Trailer     *string           `json:"trailer"`
	Platforms   []string          `json:"platforms"`
	Tags        []string          `json:"tags"`
	Images      map[string]string `json:"images"`
	Description string            `json:"description"`
}

type ContactMessage struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

var games []Game

func main() {
	port := getEnv("PORT", "8080")
	allowOrigin := getEnv("ALLOW_ORIGIN", "*")

	// Load games from JSON file
	if err := loadGames("data/games.json"); err != nil {
		log.Fatalf("Failed to load games: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, allowOrigin)
		if r.Method == http.MethodOptions {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	mux.HandleFunc("/api/games", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, allowOrigin)
		if r.Method == http.MethodOptions {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	})

	mux.HandleFunc("/api/contact", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, allowOrigin)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var msg ContactMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		log.Printf("Contact message received: %+v\n", msg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resp := map[string]interface{}{
			"received": true,
			"at":       time.Now().UTC().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(resp)
	})

	handler := withLogging(mux)
	log.Printf("Serving API on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func loadGames(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &games); err != nil {
		return err
	}
	return nil
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func setCORS(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
