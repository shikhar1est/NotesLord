package database

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"noteslord/models"
)

var DB *gorm.DB 

func ConnectDB(){
	var err error
	DB,err=gorm.Open(sqlite.Open("notes.db"),&gorm.Config{}) //We establish the connection and
	// "sqlite.Open("notes.db")" specifies specifies to use the SQLite driver and connect to a file named notes.db .
	//  If the file doesn't exist, SQLite will create it.
	//gorm.Open(...) attempts to open the connection using the 
	// specified driver and passes an optional gorm.Config{} for configuration.
	//  The result (the connection object and an error) are assigned to DB and err.
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err=DB.AutoMigrate(&models.Note{}, &models.User{}) //This is a crucial step: GORM inspects the models.Note struct
	//  and ensures the corresponding database table (notes) exists and has all the necessary columns
	//  (ID, Title, Content, etc.). If the table doesn't exist, it creates it.
	//  If the table exists but is missing a column, it adds the column.
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database connected and migrated successfully")
}