package database

import (
	"errors"
	"log"
)

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return Chirp{}, err
	}

	chirp := Chirp{}
	chirp.Body = body
	chirp.Id = int32(len(dbStructure.Chirps) + 1)

	dbStructure.Chirps[int(chirp.Id)] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return []Chirp{}, err
	}
	chirps := []Chirp{}
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetChirp(chirpID int32) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return Chirp{}, err
	}
	chirp, isFound := dbStructure.Chirps[int(chirpID)]
	if !isFound {
		return Chirp{}, errors.New("Not Found")
	}
	return chirp, nil
}
