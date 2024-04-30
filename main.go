package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ross-brown/chirpy-go/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	debugPtr := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debugPtr {
		err := resetDatabase()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("database.json wiped")
	}

	const PORT = "8080"
	const FILEPATH_ROOT = "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(FILEPATH_ROOT)))
	mux.Handle("/app/*", http.StripPrefix("/app", fsHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}

	fmt.Printf("Server listening on port %s...\n", PORT)
	server.ListenAndServe()
}

func resetDatabase() error {
	_, err := os.Stat("database.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = os.Remove("database.json")
	if err != nil {
		return err
	}

	return nil
}
