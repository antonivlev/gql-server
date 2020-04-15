/*
Package database provides functions for saving and retrieving objects from the database.
TODO: add tests
*/
package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/antonivlev/gql-server/auth"
	"github.com/antonivlev/gql-server/models"
	"github.com/graph-gophers/graphql-go"
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

	db.AutoMigrate(&models.Link{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Vote{})

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
	result := gormDB.
		Preload("PostedBy").
		Preload("Votes").
		Preload("Votes.User").
		Find(&links)
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

func getUserByToken(token string) (*models.User, error) {
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

func GetUser(ctx context.Context) (*models.User, error) {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil, errors.New("Post: no token in context")
	}
	// put user into ctx instead
	user, errUser := getUserByToken(token)
	if errUser != nil {
		return nil, errUser
	}
	return user, nil
}

func CreateVote(voterID, linkID graphql.ID) (*models.Vote, error) {
	if doesVoteExist(voterID, linkID) {
		return nil, errors.New("User already voted for this Link")
	}
	vote := models.Vote{
		UserID: voterID,
		LinkID: linkID,
	}
	result := gormDB.Create(&vote)
	if result.Error != nil {
		return nil, result.Error
	}

	var savedVote models.Vote
	res := gormDB.Preload("User").Preload("Link").Where("id = ?", vote.ID).Find(&savedVote)
	if res.Error != nil {
		fmt.Println(res.Error.Error())
	}
	return &savedVote, nil
}

func doesVoteExist(voterID, linkID graphql.ID) bool {
	// todo: implement
	return false
}
