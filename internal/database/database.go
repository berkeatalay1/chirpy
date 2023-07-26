package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

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

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	user := User{}
	user.Email = email
	user.Id = int32(len(dbStructure.Users) + 1)
	user.Password = password

	dbStructure.Users[int(user.Id)] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	return user, nil
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

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	file, err := os.ReadFile(db.path)

	dbStructure := DBStructure{}
	if err != nil {
		log.Fatal(err)
		return DBStructure{}, err
	}

	err = json.Unmarshal(file, &dbStructure)

	if err != nil {
		log.Fatal(err)
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := json.Marshal(dbStructure)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.WriteFile(db.path, data, 0666)
	return err
}
