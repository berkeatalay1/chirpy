package database

import "sync"

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Body string `json:"body"`
	Id   int32  `json:"id"`
}

type User struct {
	Email    string `json:"email"`
	Id       int32  `json:"id"`
	Password string
}
