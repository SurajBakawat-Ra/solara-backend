
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

var games = []Game{
	{
		ID:    "vector-horizon",
		Title: "Vector Horizon",
		Stores: map[string]string{
			"steam": "https://store.steampowered.com/app/3801540/Vector_Horizon/",
		},
		Trailer:   strPtr("Q_pv6QVvIjA"),
		Platforms: []string{"PC"},
		Tags:      []string{"Arcade", "Action"},
		Images: map[string]string{
			"cover": "https://img.youtube.com/vi/Q_pv6QVvIjA/hqdefault.jpg",
		},
		Description: "Fast-paced action set in a stylized vector world. Official store and trailer links below.",
	},
	{
		ID:    "gunboxing",
		Title: "GunBoxing",
		Stores: map[string]string{
			"steam": "https://store.steampowered.com/app/1978090/GunBoxing/",
		},
		Trailer:   strPtr("YiDwlVA0btQ"),
		Platforms: []string{"PC"},
		Tags:      []string{"Fighting", "Action"},
		Images: map[string]string{
			"cover": "https://img.youtube.com/vi/YiDwlVA0btQ/hqdefault.jpg",
		},
		Description: "Punch, shoot, and style—an over-the-top brawler. Official store and trailer links below.",
	},
	{
		ID:    "cosmo-war",
		Title: "Cosmo War",
		Stores: map[string]string{
			"googlePlay": "https://play.google.com/store/apps/details?id=com.solaragames.cosmowar",
		},
		Trailer:   nil,
		Platforms: []string{"Android"},
		Tags:      []string{"Arcade", "Casual"},
		Images:    map[string]string{"cover": ""},
		Description: "Mobile arcade action. Google Play link below.",
	},
	{
		ID:    "almanac",
		Title: "Almanac",
		Stores: map[string]string{
			"googlePlay": "https://play.google.com/store/apps/details?id=com.RaAten.Almanac",
		},
		Trailer:   strPtr("pW4nIEWJsmc"),
		Platforms: []string{"Android"},
		Tags:      []string{"Puzzle", "Casual"},
		Images:    map[string]string{"cover": "https://img.youtube.com/vi/pW4nIEWJsmc/hqdefault.jpg"},
		Description: "A thoughtful mobile experience. Trailer and store link below.",
	},
	{
		ID:    "chicken-bounce",
		Title: "Chicken Bounce",
		Stores: map[string]string{
			"googlePlay": "https://play.google.com/store/apps/details?id=com.solaragames.chickenbounce",
		},
		Trailer:   nil,
		Platforms: []string{"Android"},
		Tags:      []string{"Arcade", "Casual"},
		Images:    map[string]string{"cover": ""},
		Description: "Light, bouncy fun on mobile. Google Play link below.",
	},
	{
		ID:    "tappy-fly",
		Title: "Tappy Fly",
		Stores: map[string]string{
			"googlePlay": "https://play.google.com/store/apps/details?id=com.solaragames.tappyfly",
		},
		Trailer:   nil,
		Platforms: []string{"Android"},
		Tags:      []string{"Arcade", "Casual"},
		Images:    map[string]string{"cover": ""},
		Description: "Tap to fly and dodge in this pick-up-and-play mobile game.",
	},
}

func strPtr(s string) *string { return &s }

func main() {
	port := getEnv("PORT", "8080")
	allowOrigin := getEnv("ALLOW_ORIGIN", "*")

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
		// TODO: Integrate email (SMTP or a provider). For now, just echo success.
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
