package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"backend/routes"
	"backend/store"
	"backend/utils"

	"github.com/go-sql-driver/mysql"
)

func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}

func loadEnvFiles() {
	for _, path := range []string{".env", "../.env"} {
		loadEnvFile(path)
	}
}

func ensureSchema() {
	_, err := store.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(16) PRIMARY KEY,
			secret_code VARCHAR(32) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			is_admin BOOLEAN DEFAULT FALSE
		)
	`)
	if err != nil {
		log.Printf("Schema init warning for users: %v", err)
	}

	_, err = store.DB.Exec(`
		CREATE TABLE IF NOT EXISTS complaints (
			id VARCHAR(16) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			summary TEXT NOT NULL,
			severity INT NOT NULL,
			resolved BOOLEAN DEFAULT FALSE,
			user_id VARCHAR(16) NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Printf("Schema init warning for complaints: %v", err)
	}
}

func main() {
	loadEnvFiles()

	// Connect to MySQL
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "complaint_db"
	}

cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   net.JoinHostPort(host, port),
		DBName: name,
		TLSConfig: "skip-verify",
	}

	if err := store.InitDB(cfg.FormatDSN()); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ensureSchema()

	// Create Default Admin if it doesn't exist
	var adminSecret string
	err := store.DB.QueryRow("SELECT secret_code FROM users WHERE email = ? AND is_admin = true LIMIT 1", "admin@example.com").Scan(&adminSecret)
	if err != nil {
		adminID, _ := utils.GenRandomHex(8)
		adminSecret, _ = utils.GenRandomHex(16)

		_, err := store.DB.Exec("INSERT INTO users (id, secret_code, name, email, is_admin) VALUES (?, ?, 'admin', 'admin@example.com', true)", adminID, adminSecret)
		if err != nil {
			log.Printf("Failed to create default admin: %v", err)
		} else {
			log.Println("Default admin secret:", adminSecret)
		}
	} else {
		log.Println("Default admin secret:", adminSecret)
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080"
	}

	log.Println("Server running on", serverPort)

	routes.SetupRoutes()

	log.Fatal(http.ListenAndServe(serverPort, nil))
}
