package config

import (
	"log"
	"os"
)

type Config struct {
	DetabaseURL string
	JWTSecret   string
	Port        string 
}

