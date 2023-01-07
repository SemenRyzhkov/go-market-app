package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetServerAddress() string {
	return os.Getenv("RUN_ADDRESS")
}

func GetKey() string {
	return os.Getenv("SECRET_KEY")
}

func GetDBAddress() string {
	return os.Getenv("DATABASE_URI")
}

func GetAccrualAddress() string {
	return os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
}

func LoadEnvironments(envFilePath string) {
	err := godotenv.Load(envFilePath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
