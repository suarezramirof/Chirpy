package database

import (
	"errors"
)

func (db *DB) AddUser(email string, password string) (ReturnedUser, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return ReturnedUser{}, err
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
		return ReturnedUser{}, err
	}
	return ReturnedUser{Email: usr.Email, Id: usr.Id}, nil
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

func (db *DB) UpdateUser(id int, email string, password string) (ReturnedUser, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return ReturnedUser{}, err
	}
	usr, ok := dbStructure.Users[id]
	if !ok {
		return ReturnedUser{}, errors.New("user not found")
	}
	usr.Email = email
	usr.Password = password
	dbStructure.Users[id] = usr
	err = db.writeDB(dbStructure)
	if err != nil {
		return ReturnedUser{}, err
	}
	updatedUser := ReturnedUser{Email: usr.Email, Id: usr.Id}
	return updatedUser, nil
}