package models

import (
	"log"
	"time"

	"github.com/mizhexiaoxiao/k8s-api-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Model struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deleteAt" gorm:"index"`
}

// Setup initializes the database instance
func Setup() {
	log.Println("Setting up database connection")
	var err error

	dsn := config.DBdsn()

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	DB.AutoMigrate(&ClusterModel{})
}
