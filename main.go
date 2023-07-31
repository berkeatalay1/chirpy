package main

import (
	"log"
	"net/http"
	"os"

	"github.com/berkeatalay/chirpy/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	const port = "8080"
	godotenv.Load()
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", readyHandler)
	apiRouter.Post("/chirps", apiCfg.add_chirp)
	apiRouter.Get("/chirps", apiCfg.get_chirps)
	apiRouter.Get("/chirps/{CHIRPID}", apiCfg.get_chirp)
	apiRouter.Post("/users", apiCfg.add_user)
	apiRouter.Put("/users", apiCfg.update_user)
	apiRouter.Post("/login", apiCfg.login)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.metricsHandler)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	r.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)
	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
