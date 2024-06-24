package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int64
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
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsReader)
	mux.HandleFunc("POST /api/validate_chirp", chirpHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetMetrics)
	log.Printf("Listening on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}