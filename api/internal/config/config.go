package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		Port:               port,
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:    os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
	}

	if cfg.SupabaseURL == "" {
		log.Fatal("SUPABASE_URL is required")
	}
	if cfg.SupabaseAnonKey == "" {
		log.Fatal("SUPABASE_ANON_KEY is required")
	}

	return cfg
}
