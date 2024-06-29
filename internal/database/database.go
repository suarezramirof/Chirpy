package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
	"github.com/suarezramirof/Chirpy/shared"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type RefreshToken struct {
	UserId    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type DBStructure struct {
	Chirps        map[int]shared.Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

// type Chirp struct {
// 	Body string `json:"body"`
// 	Id   int    `json:"id"`
// 	AuthorId int `json:"author_id"`
// }

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Id       int    `json:"id"`
	IsChirpyRed bool `json:"is_chirpy_red"`
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
		Chirps:        make(map[int]shared.Chirp),
		Users:         make(map[int]User),
		RefreshTokens: make(map[string]RefreshToken),
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
