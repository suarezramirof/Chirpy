package main

import (
	"log"
	"net/http"
	"strconv"
)

type apiConfig struct {
	fileserverHits int64
}

var apiCfg = apiConfig{
	fileserverHits: 0,
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", metricsHandler)
	mux.HandleFunc("/reset", resetHandler)
	log.Printf("Listening on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + apiCfg.getHits()))
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	apiCfg.resetHits()
}



func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		cfg.fileserverHits++
	})
}

func (cfg *apiConfig) getHits() string {
	return strconv.Itoa(int(cfg.fileserverHits))
}

func (cfg *apiConfig) resetHits() {
	cfg.fileserverHits = 0
}