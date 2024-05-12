package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Phone struct {
	gorm.Model
	UserID    uint
	Prefix    string
	Number    uint
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}
type User struct {
	gorm.Model
	FirstName string `gorm:"uniqueIndex"`
	LastName  string `gorm:"uniqueIndex"`
	Email     string `gorm:"not null"`
	Country   string
	Role      string `gorm:"default:user;not null"`
	Age       uint   `gorm:"not null;size:3"`
	PhoneID   uint
	Phone     Phone
	Phones    []Phone
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

type FaqType struct {
	gorm.Model
	Name        string
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
}

type Faq struct {
	gorm.Model
	ParentID  uint
	FaqTypeID uint
	FaqType   FaqType
	Question  string
	Answer    string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func StartDB() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Kuala_Lumpur"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	err = db.AutoMigrate(&User{}, &Phone{}, &FaqType{}, &Faq{})
	if err != nil {
		panic("failed to perform migrations: " + err.Error())
	}
}
