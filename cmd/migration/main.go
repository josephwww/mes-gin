package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"mes-gin/models"
	"mes-gin/utils"
)

func main() {
	dsn := "host=localhost user=mes-gin-user password=mysecretpassword dbname=mes_gin_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// List of models to drop
	modelsToDrop := []interface{}{
		&models.User{},
		&models.Organization{},
		&models.Role{},
		// Add more models here as needed
	}

	// Loop over each model and drop the table if it exists
	for _, model := range modelsToDrop {
		if db.Migrator().HasTable(model) {
			err := db.Migrator().DropTable(model)
			if err != nil {
				log.Fatal("failed to drop table for model:", err)
			}
		}
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Organization{}, &models.User{}, &models.Role{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}

	log.Println("Super user created successfully")

	dsn = "host=localhost user=mes-gin-user password=mysecretpassword dbname=mes_gin_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Organization{}, &models.User{}, &models.Role{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}

	pwd, _ := utils.HashPassword("12345678")
	db.Create(&models.User{
		Name:     "super_user",
		Phone:    "18888888888",
		Password: pwd,
	})

}
