package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (cfg *apiConfig) get_chirp(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := chi.URLParam(r, "CHIRPID")
	chirpId, err := strconv.Atoi(chirpIdStr)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	chirp, err := cfg.DB.GetChirp(int32(chirpId))
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	respondWithJSON(w, 200, chirp)
}
