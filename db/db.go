package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string    `gorm:"uniqueIndex"`
	LastName  string    `gorm:"uniqueIndex"`
	Email     string    `gorm:"not null"`
	Country   string    `gorm:"not null"`
	Role      string    `gorm:"not null"`
	Age       int       `gorm:"not null;size:3"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func StartDB() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Kuala_Lumpur"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to perform migrations: " + err.Error())
	}
}
