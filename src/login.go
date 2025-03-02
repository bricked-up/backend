package backend 

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
	"golang.org/x/crypto/bcrypt"   // Library for secure password hashing
)

// Credentials represents the user's login credentials (email and password).
type Credentials struct {
	Email    string `json:"email"`    // User's email address
	Password string `json:"password"` // User's password (plaintext during login)
}

// Session represents a user session with an expiration time and a unique token.
type Session struct {
	UserID    int       `json:"user_id"`    // ID of the user associated with the session
	ExpiresAt time.Time `json:"expires_at"` // Time when the session expires
	Token     string    `json:"token"`      // Unique token for the session
}

// db is a global variable that holds the database connection.
var db *sql.DB

// init is called when the program starts. It initializes the database connection.
func init() {
	var err error
	// Open a connection to the SQLite database file.
	db, err = sql.Open("sqlite3", "./bricked-up_prod.db")
	if err != nil {
		// If the connection fails, log the error and stop the program.
		log.Fatal("Database connection failed:", err)
	}
}

// loginHandler handles HTTP POST requests to the /login endpoint for user authentication.
func login(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid request method: %s", r.Method)
	}
	return nil
}

	// Decode the JSON request body into a Credentials struct.
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Query the database to retrieve the user's ID and hashed password.
	var userID int
	var storedPassword string
	err := db.QueryRow("SELECT id, password FROM USER WHERE email = ?", creds.Email).Scan(&userID, &storedPassword)
	if err != nil {
		// If the query fails (e.g., user not found), log the error and return an unauthorized response.
		log.Println("Database error:", err)
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password.
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		// If the passwords don't match, return an unauthorized response.
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	// Generate a secure session token for the user.
	token, err := generateToken()
	if err != nil {
		// If token generation fails, log the error and return an internal server error.
		log.Println("Token generation error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the session expiration time to 24 hours from now.
	expiration := time.Now().Add(24 * time.Hour)
	// Insert the new session into the database.
	_, err = db.Exec("INSERT INTO SESSION (userid, expires, token) VALUES (?, ?, ?)", userID, expiration, token)
	if err != nil {
		// If session creation fails, log the error and return an internal server error.
		log.Println("Session creation error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create a Session struct to return to the client.
	session := Session{
		UserID:    userID,
		ExpiresAt: expiration,
		Token:     token,
	}

	// Set the response content type to JSON and encode the session as JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// generateToken generates a secure, unique token for user sessions.
func generateToken() (string, error) {
	// Placeholder for secure token generation.
	// In a real application, use a cryptographically secure random generator.
	return "secure-random-token", nil
}