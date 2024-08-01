package db

import (
        "gorm.io/driver/sqlite"
        "gorm.io/gorm"
		
)

var DB *gorm.DB

type Task struct {
	gorm.Model
	Title   string
	Content string
	Complete bool
}

func ConnectDatabase() {
        var err error
        DB, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
        if err != nil {
                panic("failed to connect database")
        }

        DB.AutoMigrate(&Task{})
}


