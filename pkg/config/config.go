package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DbConfig  *DBConfig
	JwtConfig *JWTConfig
}

type DBConfig struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}
type JWTConfig struct {
	SecretKey string
	Prefix    string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, relying on environment variables")
	}

	return &Config{
		DbConfig: &DBConfig{
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
		},
		JwtConfig: &JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			Prefix:    os.Getenv("JWT_PREFIX"),
		},
	}
}
