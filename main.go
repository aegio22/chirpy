package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/aegio22/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}
	dbQueries := database.New(db)
	cfg := apiConfig{db: dbQueries}
	multiPlexer := http.NewServeMux()
	const rootDir = http.Dir(".")
	//initialize fileserver with hits counter
	multiPlexer.Handle("/app/",
		http.StripPrefix("/app/",
			cfg.middlewareMetricsInc(http.FileServer(rootDir))))

	//initialize core endpoints/handlers
	multiPlexer.HandleFunc("GET /api/healthz", handlerReadiness)
	multiPlexer.HandleFunc("GET /admin/metrics", cfg.handlerCountReqs)
	multiPlexer.HandleFunc("POST /admin/reset", cfg.handlerResetHitCount)
	multiPlexer.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	//initialize and start server
	server := http.Server{
		Addr:    ":8080",
		Handler: multiPlexer,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}

}

func (cfg *apiConfig) handlerCountReqs(rw http.ResponseWriter, req *http.Request) {
	hits := cfg.fileserverHits.Load()
	hitsReport := fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", hits)

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Write([]byte(hitsReport))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(rw, req)
	})

}
