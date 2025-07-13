package Storage

import (
	"fmt"
	"github.com/coderconquerer/go-login-app/common"
	"github.com/coderconquerer/go-login-app/configs"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConnection struct {
	conn *gorm.DB
}

func GetMySQLConnection(cfg *configs.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, common.NewDatabaseError(err)
	}

	return db, nil
}

func GetNewMySQLConnection(db *gorm.DB) *MySQLConnection {
	return &MySQLConnection{conn: db}
}
