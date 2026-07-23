package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port               string
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
}

type yamlConfig struct {
	Supabase struct {
		URL        string `yaml:"url"`
		AnonKey    string `yaml:"anon_key"`
		ServiceKey string `yaml:"service_key"`
	} `yaml:"supabase"`
	Server struct {
		GoPort int `yaml:"go_port"`
	} `yaml:"server"`
}

func resolveEnv(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		envName := value[2 : len(value)-1]
		return os.Getenv(envName)
	}
	return value
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	yamlPath := os.Getenv("CONFIG_YAML")
	if yamlPath == "" {
		yamlPath = "../config.yaml"
	}

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("Cannot read %s: %v", yamlPath, err)
	}

	var yc yamlConfig
	if err := yaml.Unmarshal(data, &yc); err != nil {
		log.Fatalf("Cannot parse %s: %v", yamlPath, err)
	}

	port := os.Getenv("PORT")
	if port == "" && yc.Server.GoPort != 0 {
		port = fmt.Sprintf("%d", yc.Server.GoPort)
	}
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		Port:               port,
		SupabaseURL:        yc.Supabase.URL,
		SupabaseAnonKey:    resolveEnv(yc.Supabase.AnonKey),
		SupabaseServiceKey: resolveEnv(yc.Supabase.ServiceKey),
	}

	if cfg.SupabaseURL == "" {
		log.Fatal("supabase.url is required in config.yaml")
	}
	if cfg.SupabaseAnonKey == "" {
		log.Fatal("supabase.anon_key is required in config.yaml")
	}

	return cfg
}
