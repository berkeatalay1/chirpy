package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type request struct {
	Body string `json:"body"`
}

type response struct {
	Valid bool `json:"valid"`
}

type error struct {
	Error string `json:"error"`
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		error := error{Error: "Something went wrong"}
		dat, err := json.Marshal(error)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	if len(req.Body) > 140 {
		error := error{Error: "Chirp is too long"}
		dat, err := json.Marshal(error)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	respBody := response{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
