package helpers

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupConnection(dbDns string) {
	connection, err := gorm.Open(mysql.Open(dbDns), &gorm.Config{})
	if err != nil {
		panic("could not connect to the database")
	}

	DB = connection

}
