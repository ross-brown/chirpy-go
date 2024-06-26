package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	IsChirpyRed    bool   `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		fmt.Println("ERROR LOADING DB")
		return User{}, err
	}

	_, err = db.GetUserByEmail(email)
	if err == nil {
		fmt.Println("USER EXISTS ALREADY")
		return User{}, errors.New("this email already exists")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println("ERROR GETTING HASH")
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:             id,
		Email:          email,
		HashedPassword: string(hashedBytes),
		IsChirpyRed:    false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		fmt.Println("ERROR WRITING TO DB")
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByID(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, dbUser := range dbStructure.Users {
		if dbUser.Email == email {
			return dbUser, nil
		}
	}

	return User{}, errors.New("resource does not exist")
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("resource does not exist")
	}

	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeChirpyRed(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {

		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {

		return User{}, errors.New("resource does not exist")
	}

	user.IsChirpyRed = true
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
