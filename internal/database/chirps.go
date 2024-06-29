package database

import (
	"errors"
	s "github.com/suarezramirof/Chirpy/shared"
)

func (db *DB) CreateChirp(body string, userId int) (s.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return s.Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := s.Chirp{
		Body:     body,
		Id:       id,
		AuthorId: userId,
	}
	dbStructure.Chirps[id] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return s.Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) GetChirps() ([]s.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]s.Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetChirp(id int) (s.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return s.Chirp{}, err
	}
	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return s.Chirp{}, errors.New("chirp not found")
	}
	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	_, ok := dbStructure.Chirps[id]
	if !ok {
		return errors.New("chirp not found")
	}
	delete(dbStructure.Chirps, id)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
