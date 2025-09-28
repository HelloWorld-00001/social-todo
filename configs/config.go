package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
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

// InitEnv loads environment variables from .env file
// If envPath is provided, it uses that file, otherwise falls back to default.
// InitEnv looks for a .env file by walking upward from the current directory
// until it finds one, and loads it into environment variables.
func InitEnv() {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Could not get working directory: %v", err)
		return
	}

	// Walk upward looking for .env
	var envPath string
	for dir := wd; dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		try := filepath.Join(dir, ".env")
		if _, err := os.Stat(try); err == nil {
			envPath = try
			break
		}
	}

	if envPath == "" {
		log.Println("No .env file found, relying on system environment variables")
		return
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Failed to load .env file at %s: %v", envPath, err)
	} else {
		log.Printf("Loaded environment variables from %s", envPath)
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
