package database

import (
	"errors"
	"time"
)

func (db *DB) SaveRefreshToken(userId int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	refreshToken := RefreshToken{
		UserId:    userId,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	dbStructure.RefreshTokens[token] = refreshToken
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CheckRefreshToken(token string) (int, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return 0, err
	}
	refTok, ok := dbStructure.RefreshTokens[token]
	// if token not found or expired
	if !ok || refTok.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("invalid or expired refresh token")
	}
	userId := refTok.UserId
	return userId, nil
}

func (db *DB) DeleteRefreshToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	delete(dbStructure.RefreshTokens, token)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
