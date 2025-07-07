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

type AWSConfig struct {
	SecretKey string
	PublicKey string
	Bucket    string
	Region    string
	Domain    string
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

func LoadAWSConfig() *AWSConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, relying on environment variables")
	}

	return &AWSConfig{
		PublicKey: os.Getenv("AWS_PUBLIC_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region:    os.Getenv("AWS_REGION"),
		Domain:    os.Getenv("AWS_DOMAIN"),
		Bucket:    os.Getenv("AWS_BUCKET"),
	}
}
