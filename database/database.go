/*
Package database provides functions for saving and retrieving objects from the database.
TODO: add tests
*/
package database

import (
	"errors"

	"github.com/antonivlev/gql-server/auth"
	"github.com/antonivlev/gql-server/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	gormDB *gorm.DB
)

// Connects the database
func SetupDatabase() error {
	db, errConnect := gorm.Open("postgres",
		"host=localhost"+
			" port=5432"+
			" user=pathfinder"+
			" dbname=hackernews"+
			" password=pathfinder"+
			" sslmode=disable",
	)
	if errConnect != nil {
		return errConnect
	}

	// todo: check for migration error
	db.AutoMigrate(&models.Link{})
	db.AutoMigrate(&models.User{})

	gormDB = db
	return nil
}

func CreateLink(link models.Link) (*models.Link, error) {
	result := gormDB.Create(&link)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func GetAllLinks() ([]models.Link, error) {
	links := []models.Link{}
	result := gormDB.Preload("PostedBy").Find(&links)
	if result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}

func CreateUser(email, password, name string) (*models.User, error) {
	if doesUserWithEmailExist(email) {
		return nil, errors.New("User with email: " + email + " already exists")
	}

	hashedPasswordBytes, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHash != nil {
		return nil, errHash
	}

	newUser := models.User{
		Email:    email,
		Password: string(hashedPasswordBytes),
		Name:     name,
	}
	result := gormDB.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser, nil
}

func doesUserWithEmailExist(email string) bool {
	user := models.User{}
	gormDB.Where("email = ?", email).First(&user)
	return user.ID != ""
}

func GetUserByCredentials(email, password string) (*models.User, error) {
	user := models.User{}
	gormDB.Where("email = ?", email).First(&user)
	if user.ID == "" {
		return nil, errors.New("no user with email " + email)
	}

	errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errCompare != nil {
		return nil, errCompare
	} else {
		return &user, nil
	}
}

func GetUserByToken(token string) (*models.User, error) {
	userID := auth.GetUserIDFromToken(token)
	if userID == "" {
		return nil, errors.New("Could not determine user from token")
	}

	var user models.User
	gormDB.Where("id = ?", userID).Preload("Links").Find(&user)
	if user.ID == "" {
		return nil, errors.New("No user with decoded ID from token")
	}

	return &user, nil
}
