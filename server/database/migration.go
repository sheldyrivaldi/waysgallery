package database

import (
	"fmt"
	"waysgallery/models"
	"waysgallery/pkg/mysql"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Photo{},
		&models.Order{},
		&models.Bank{},
		&models.Withdrawal{},
		&models.PhotoProject{},
		&models.Project{},
	)

	if err != nil {
		fmt.Println(err.Error())
		panic("Migration failed!")
	}

	fmt.Println("Migration successfully completed!")
}
