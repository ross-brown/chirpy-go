package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const PORT = "8080"
	const FILEPATH_ROOT = "."

	apiCfg := &apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(FILEPATH_ROOT)))
	mux.Handle("/app/*", http.StripPrefix("/app", fsHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerResetMetrics)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}

	fmt.Printf("Server listening on port %s...\n", PORT)
	server.ListenAndServe()
}
