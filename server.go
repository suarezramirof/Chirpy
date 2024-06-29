package main

import (
	"log"
	"net/http"
	"os"
	database "github.com/suarezramirof/Chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB *database.DB
	jwtSecret string
	polkaKey string
}

func main() {
	godotenv.Load()
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
		jwtSecret: os.Getenv("JWT_SECRET"),
		polkaKey: os.Getenv("POLKA_KEY"),
	}
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg.DB = db
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsReader)
	mux.HandleFunc("POST /api/chirps", apiCfg.chirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/login", apiCfg.loginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeToken)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeUser)
	mux.HandleFunc("GET /api/chirps", apiCfg.chirpsGetter)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.deleteChirp)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.chirpGetter)
	mux.HandleFunc("/api/reset", apiCfg.resetMetrics)
	log.Printf("Listening on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}