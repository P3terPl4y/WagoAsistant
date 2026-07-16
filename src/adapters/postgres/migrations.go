package postgres

import (
	"App/src/pkg/logger"
	"database/sql"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Connect opens a database connection and initializes the schema.
func Connect(databaseURL string, log logger.Logger) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Database ping failed")
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	runMigrations(db, log)
	log.Info().Msg("PostgreSQL database initialized")
	return db
}

// runMigrations creates tables and applies schema migrations.
func runMigrations(db *sql.DB, log logger.Logger) {
	createTables := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		phone TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS bots (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		blocked BOOLEAN DEFAULT FALSE,
		session_file TEXT,
		payment_status TEXT DEFAULT 'free',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS prompts (
		bot_id INTEGER PRIMARY KEY REFERENCES bots(id) ON DELETE CASCADE,
		prompt TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS subscriptions (
		bot_id INTEGER PRIMARY KEY REFERENCES bots(id) ON DELETE CASCADE,
		tier TEXT NOT NULL DEFAULT 'free',
		msg_limit INTEGER NOT NULL DEFAULT 10,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL
	);
	CREATE TABLE IF NOT EXISTS chat_history (
		id SERIAL PRIMARY KEY,
		bot_id INTEGER NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
		user_jid TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS oauth_tokens (
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		provider TEXT NOT NULL,
		refresh_token TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, provider)
	);`

	if _, err := db.Exec(createTables); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	_, _ = db.Exec(`ALTER TABLE bots ADD COLUMN IF NOT EXISTS payment_status TEXT DEFAULT 'free'`)
	_, _ = db.Exec(`ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS tier TEXT DEFAULT 'free'`)
	_, _ = db.Exec(`ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS msg_limit INTEGER DEFAULT 10`)
}

// EnsureAdmin creates the default admin user if none exists.
func EnsureAdmin(db *sql.DB, username, email, phone, password string, log logger.Logger) {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'admin'`).Scan(&count); err != nil {
		log.Fatal().Err(err).Msg("Failed to query admin count")
	}
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		_, err := db.Exec(
			`INSERT INTO users (username, email, phone, password_hash, role) VALUES ($1, $2, $3, $4, 'admin')`,
			username, email, phone, string(hashed),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create admin user")
		}
		log.Info().Str("username", username).Msg("Admin user created")
	}
}
