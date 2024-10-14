package core

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Plugin struct {
	gorm.Model
	Name string
	Url  string
}

var Database *gorm.DB

func GetDatabase() *gorm.DB {
	return Database
}

func LoadDatabase() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Plugin{})

	Database = db
}
