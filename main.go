package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", readyHandler)
	apiRouter.Post("/validate_chirp", validate_chirp)

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
