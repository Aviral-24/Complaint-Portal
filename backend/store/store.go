package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

const MASTER_ADMIN_KEY = "CHANGE_ME_MASTER_KEY_123" //[cite: 8]

// User struct kept matching the original JSON format[cite: 8]
type User struct {
	ID         string   `json:"id"`
	SecretCode string   `json:"secret_code"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Complaints []string `json:"complaints"`
	IsAdmin    bool     `json:"is_admin"`
}

// Complaint struct kept matching the original JSON format[cite: 8]
type Complaint struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Severity int    `json:"severity"`
	Resolved bool   `json:"resolved"`
	UserID   string `json:"user_id"`
}

func InitDB(dsn string) error {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return DB.Ping()
}

func GetUserBySecret(secret string) (*User, error) {
	u := &User{}
	err := DB.QueryRow("SELECT id, secret_code, name, email, is_admin FROM users WHERE secret_code = ?", secret).Scan(&u.ID, &u.SecretCode, &u.Name, &u.Email, &u.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("invalid secret code") //[cite: 8]
	}

	u.Complaints = []string{}
	rows, err := DB.Query("SELECT id FROM complaints WHERE user_id = ?", u.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cid string
			rows.Scan(&cid)
			u.Complaints = append(u.Complaints, cid)
		}
	}
	return u, nil
}
