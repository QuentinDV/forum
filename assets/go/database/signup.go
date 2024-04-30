package database

// Importing necessary packages
import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SignUpError struct holds the error information for sign up process.
type SignUpError struct {
	EmailError    bool
	UsernameError bool
}

// CreateAccount function creates a new account in the database.
// It takes email, password, username and admin status as input.
// It returns the created account, any sign up error and error if any.
func CreateAccount(email, password, username string, isAdmin bool) (Account, SignUpError, error) {
	var account Account
	// Connexion à la base de données
	db, err := ConnectDB("database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the last ID from the database
	var lastID string
	row := db.QueryRow("SELECT id FROM accounts ORDER BY id DESC LIMIT 1")
	err = row.Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			lastID = "0"
		} else {
			return account, SignUpError{}, err
		}
	}

	// Increment the last ID
	newID := incrementID(lastID)

	// Vérifier si l'email est déjà pris
	emailTaken, err := IsEmailTaken(db, email)
	if err != nil {
		return account, SignUpError{}, err
	}

	// Vérifier si le pseudonyme est déjà pris
	usernameTaken, err := IsUsernameTaken(db, username)
	if err != nil {
		return account, SignUpError{}, err
	}

	if emailTaken && usernameTaken {
		return account, SignUpError{EmailError: true, UsernameError: true}, nil
	}
	if emailTaken {
		return account, SignUpError{EmailError: true, UsernameError: false}, nil
	} else if usernameTaken {
		return account, SignUpError{EmailError: false, UsernameError: true}, nil
	}

	// Exemple d'utilisation : Création et insertion d'un nouveau compte
	newAccount := Account{
		Id:           newID,
		Email:        email,
		Password:     hashPasswordSHA256(password),
		Username:     username,
		ImageUrl:     "https://i.pinimg.com/474x/63/bc/94/63bc9469cae29b897565a08f0647db3c.jpg",
		IsAdmin:      isAdmin,
		IsBan:        false,
		CreationDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = InsertAccount(db, newAccount)
	if err != nil {
		log.Fatal(err)
	}

	return newAccount, SignUpError{}, nil
}

// IsEmailTaken function checks if an email is already taken in the database.
// It takes a database connection and an email as input.
// It returns a boolean indicating if the email is taken and an error if any.
func IsEmailTaken(db *sql.DB, email string) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM accounts WHERE email = ?", email)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsUsernameTaken function checks if a username is already taken in the database.
// It takes a database connection and a username as input.
// It returns a boolean indicating if the username is taken and an error if any.
func IsUsernameTaken(db *sql.DB, username string) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM accounts WHERE username = ?", username)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// incrementID function increments the last ID.
// It takes the last ID as input and returns the incremented ID.
func incrementID(lastID string) string {
	var id int
	_, err := fmt.Sscanf(lastID, "%d", &id)
	if err != nil {
		log.Fatal(err)
	}
	id++
	return fmt.Sprintf("%d", id)
}

// hashPasswordSHA256 function hashes a password using SHA256 encryption.
// It takes a password as input and returns the hashed password.
func hashPasswordSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
