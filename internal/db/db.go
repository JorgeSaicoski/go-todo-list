package db

import (
        "gorm.io/driver/sqlite"
        "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
        var err error
        DB, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
        if err != nil {
                panic("failed to connect database")
        }

        DB.AutoMigrate(&models.Task{})
}
