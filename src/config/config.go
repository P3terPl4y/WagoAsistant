package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Environment string // "development" or "production"
	ServerAddr  string

	DatabaseURL   string
	RedisURL      string // optional
	EncryptionKey []byte

	// Admin defaults
	AdminUsername string
	AdminEmail   string
	AdminPhone   string
	AdminPass    string

	// OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// AI providers
	AI AIConfig

	// Limits
	MaxBots             int
	MaxConnectRetries   int
	SubscriptionDuration time.Duration
	MaxHistory          int64
	MaxHistoryChars     int
	SessionExpiration   time.Duration
	RateLimitPerMinute  int
	MaxPromptLength     int
	MaxMsgLength        int
	AITimeoutTotal      time.Duration
	DedupWindow         time.Duration

	// Cookies
	CookieSecure bool
}

// AIConfig holds configuration for AI providers.
type AIConfig struct {
	OpenRouterKey  string
	OpenRouterURL  string
	LegacyKey      string
	LegacyURL      string
	LocalURL       string
	LocalEnabled   bool
	FreeModels     []string
}

// Load reads configuration from environment variables and returns a validated Config.
// It panics on missing required values to fail fast at startup.
func Load() *Config {
	cfg := &Config{
		Environment: getEnv("APP_ENV", "development"),
		ServerAddr:  getEnv("SERVER_ADDR", "127.0.0.1:3000"),

		DatabaseURL: requireEnv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),

		// Admin
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminEmail:   getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPhone:   getEnv("ADMIN_PHONE", "+1234567890"),
		AdminPass:    getEnv("ADMIN_PASSWORD", "admin123"),

		// OAuth
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),

		// AI
		AI: AIConfig{
			OpenRouterKey:  getEnv("OPENROUTER_API_KEY", ""),
			OpenRouterURL:  getEnv("OPENROUTER_URL", "https://openrouter.ai/api/v1/chat/completions"),
			LegacyKey:      getEnv("LEGACY_API_KEY", ""),
			LegacyURL:      getEnv("LEGACY_URL", "https://apifreellm.com/api/v1/chat"),
			LocalURL:       getEnv("LOCAL_AI_URL", "http://localhost:8080/v1/chat/completions"),
			LocalEnabled:   os.Getenv("LOCAL_AI_ENABLED") == "true",
			FreeModels:     []string{"openrouter/free"},
		},

		// Limits
		MaxBots:             50,
		MaxConnectRetries:   5,
		SubscriptionDuration: 7 * 24 * time.Hour,
		MaxHistory:          8,
		MaxHistoryChars:     2000,
		SessionExpiration:   1 * time.Hour,
		RateLimitPerMinute:  10,
		MaxPromptLength:     2000,
		MaxMsgLength:        500,
		AITimeoutTotal:      40 * time.Second,
		DedupWindow:         3 * time.Second,

		// Cookies
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
	}

	// Parse encryption key
	keyHex := requireEnv("ENCRYPTION_KEY")
	key, err := hex.DecodeString(keyHex)
	if err != nil || len(key) != 32 {
		panic("ENCRYPTION_KEY must be a 64-character hex string (32 bytes)")
	}
	cfg.EncryptionKey = key

	return cfg
}

// getEnv returns the environment variable value or a default.
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// requireEnv returns the environment variable value or panics if not set.
func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return v
}
