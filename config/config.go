package config

import (
	"fmt"
	"os"
	
	"github.com/joho/godotenv"
)


type Config struct{
	DatabaseURL string
	BaseURL string 
} 

func Load() (*Config , error){
	if err := godotenv.Load(); err != nil{
		return nil,fmt.Errorf("loading .env file: %w", err)
	}

	databaseURL := os.Getenv("Database_URL")
	if databaseURL == "" {
		return nil,fmt.Errorf("DATABASE_URL is not set")
		
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		return  nil, fmt.Errorf("BASE_URL not set")
	}

	return &Config{
		DatabaseURL: databaseURL,
		BaseURL: baseURL,
	},nil

}

