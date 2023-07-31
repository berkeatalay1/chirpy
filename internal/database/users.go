package database

import (
	"errors"
	"log"
)

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) CreateUser(email string, password string) (User, error) {
	_, isFound := db.GetUserWithEmail(email)
	if isFound {
		return User{}, errors.New("User Email Allready Exist")
	}

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

func (db *DB) GetUserWithEmail(email string) (User, bool) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
		return User{}, false
	}

	users := dbStructure.Users
	for _, user := range users {
		if user.Email == email {
			return user, true
		}
	}

	return User{}, false
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("Not Exist")
	}

	user.Email = email
	user.Password = hashedPassword
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
