package Storage

import (
	"fmt"
	"github.com/coderconquerer/go-login-app/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConnection struct {
	conn *gorm.DB
}

func GetMySQLConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetNewMySQLConnection(db *gorm.DB) *MySQLConnection {
	return &MySQLConnection{conn: db}
}
