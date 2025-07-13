package configs

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
	JwtPrefix string
}

type AWSConfig struct {
	S3Prefix  string
	SecretKey string
	PublicKey string
	Bucket    string
	Region    string
	Domain    string
}

func InitEnv() {
	err := godotenv.Load("../../.env") // adjust path relative to cmd/server/main.go
	if err != nil {
		log.Println(".env file not found, relying on environment variables")
	}
}

func Load() *Config {
	InitEnv()
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
			JwtPrefix: os.Getenv("JWT_PREFIX"),
		},
	}
}

func LoadAWSConfig() *AWSConfig {
	InitEnv()
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, relying on environment variables")
	}

	return &AWSConfig{
		S3Prefix:  os.Getenv("AWS_S3_PREFIX"),
		PublicKey: os.Getenv("AWS_S3_PUBLIC_KEY"),
		SecretKey: os.Getenv("AWS_S3_SECRET_KEY"),
		Region:    os.Getenv("AWS_S3_REGION"),
		Domain:    os.Getenv("AWS_S3_DOMAIN"),
		Bucket:    os.Getenv("AWS_S3_BUCKET"),
	}
}
