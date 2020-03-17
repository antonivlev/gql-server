/*
Package database provides functions for saving and retrieving objects from the database.
TODO: add tests
*/
package database

import (
	"github.com/antonivlev/gql-server/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	result := gormDB.Find(&links)
	if result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}
