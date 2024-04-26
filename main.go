package main

import (
	"fmt"
	"net/http"
)

func main() {
	const PORT = "8080"
	const FILEPATH_ROOT = "."

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir(FILEPATH_ROOT))))
	mux.HandleFunc("/healthz", handlerReadiness)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}

	fmt.Printf("Server listing on port %s...\n", PORT)
	server.ListenAndServe()
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
