package main

import (
	"log"
	"net/http"

	database "github.com/suarezramirof/Chirpy/internal"
)

type apiConfig struct {
	fileserverHits int
	DB *database.DB
}

func main() {
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
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
	mux.HandleFunc("GET /api/chirps", apiCfg.chirpsGetter)
	mux.HandleFunc("/api/reset", apiCfg.resetMetrics)
	log.Printf("Listening on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}