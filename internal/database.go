package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

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
	Id   int    `json:"id"`
}

type User struct {
	Email    string `json:"email"`
	Password []byte `json:"password"`
	Id       int    `json:"id"`
}

type CreatedUser struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

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
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}
	return db.writeDB(dbStructure)
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	dbStructure := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, nil
	}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		Body: body,
		Id:   id,
	}
	dbStructure.Chirps[id] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}
	return chirp, nil
}

func (db *DB) AddUser(email string, password []byte) (CreatedUser, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return CreatedUser{}, err
	}
	userId := len(dbStructure.Users) + 1
	usr := User{
		Email:    email,
		Password: password,
		Id:       userId,
	}
	dbStructure.Users[userId] = usr
	err = db.writeDB(dbStructure)
	if err != nil {
		return CreatedUser{}, err
	}
	return CreatedUser{Email: usr.Email, Id: usr.Id}, nil
}

func (db *DB) GetUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}
