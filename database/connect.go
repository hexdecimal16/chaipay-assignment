package database

import (
	"fmt"
	"strconv"

	"hexdecimal16/chaipay-assignment/config"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	configData := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_HOST"),
		port,
		config.Config("DB_NAME"),
	)

	DB, err = gorm.Open(
		"mysql",
		configData,
	)

	if err != nil {
		fmt.Println(
			err.Error(),
		)
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
}
