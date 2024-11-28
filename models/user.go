package models

import (
	"errors"
	"log"

	"abc.com/calc/db"
	"abc.com/calc/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Username string `binding:"required"`
}

type User2 struct{
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}


func GetUserById(userId int64) (string, error) {
    var username string

    query := `SELECT username FROM users WHERE id=$1`
    err := db.DB.QueryRow(query, userId).Scan(&username)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        return "", err
    }

    return username, nil
}

func GetRoleById(userId int64)(string,error){
	var role string

	query:=`SELECT role from users where id=$1`
	err:=db.DB.QueryRow(query,userId).Scan(&role)
	if err!=nil{
		log.Printf("Error fetching the role: %v",err)
		return "",err
	}
	return role,err
}

func GetEventById(eventID int64) (Event, error) {
	var event Event
	query := "SELECT id, name, description FROM events WHERE id=$1"
	err := db.DB.QueryRow(query, eventID).Scan(&event.ID, &event.Name, &event.Description)
	if err != nil {
		return event, err
	}
	return event, nil
}

func CheckUserRegistration(userID, eventID int64) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM registrations WHERE user_id=$1 AND event_id=$2"
	err := db.DB.QueryRow(query, userID, eventID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}




func (u *User) Save() error {
	log.Println("Starting Save method for User...")

	// SQL Query (role is hardcoded as 'participant')
	query := `INSERT INTO users (email, password_hash, username, role) VALUES ($1, $2, $3, 'participant') RETURNING id`

	// Hash the password
	log.Println("Hashing the password...")
	hashedPwd, err := utils.HashPassword(u.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	log.Println("Password hashed successfully.")

	// Execute the query and retrieve the auto-generated ID
	log.Printf("Executing the query with Email: %s, Username: %s", u.Email, u.Username)
	err = db.DB.QueryRow(query, u.Email, hashedPwd, u.Username).Scan(&u.ID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	log.Printf("User saved successfully with ID: %d", u.ID)

	return nil
}


func (u *User) CheckCredentials() error {
	// SQL query using PostgreSQL syntax for placeholders
	query := "SELECT id, password_hash FROM users WHERE email=$1"
	// Execute the query with the provided email
	row := db.DB.QueryRow(query, u.Email)
	var retrievedPassword string
	// Scan the result into u.ID and retrievedPassword
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		// Return a generic error to avoid leaking sensitive details
		return errors.New("invalid credentials")
	}
	// Verify the provided password with the retrieved hashed password
	validPassword := utils.CheckPassword(u.Password, retrievedPassword)
	if !validPassword {
		return errors.New("invalid credentials")
	}

	// Credentials are valid
	return nil
}

func (u *User2) CheckCredentials2() error {
	// SQL query using PostgreSQL syntax for placeholders
	query := "SELECT id, password_hash FROM users WHERE email=$1"
	// Execute the query with the provided email
	row := db.DB.QueryRow(query, u.Email)
	var retrievedPassword string
	// Scan the result into u.ID and retrievedPassword
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		// Return a generic error to avoid leaking sensitive details
		return errors.New("invalid credentials")
	}
	// Verify the provided password with the retrieved hashed password
	validPassword := utils.CheckPassword(u.Password, retrievedPassword)
	if !validPassword {
		return errors.New("invalid credentials")
	}

	// Credentials are valid
	return nil
}




