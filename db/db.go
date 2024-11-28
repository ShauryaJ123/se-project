package db

import (
	"database/sql"
	_"fmt"
	"log"
   
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func InitDB() {
	var err error
	// PostgreSQL connection string (hardcoded for simplicity)
	connStr := "user=postgres password=1234 dbname=db123 host=localhost sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	log.Println("connected")
	// Set max open and idle connections
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Create tables if not already present
	createTables()
	// insertTestData()
}

func createTables() {
	// Create USERS table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL CHECK (role IN ('participant', 'organizer', 'administrator')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Could not create USERS table: %v", err)
	}

	// Create EVENTS table
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		start_time TIMESTAMP NOT NULL,
		end_time TIMESTAMP NOT NULL,
		status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')),
		organizer_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (organizer_id) REFERENCES users (id)
	)`
	_, err = DB.Exec(createEventsTable)
	if err != nil {
		log.Fatalf("Could not create EVENTS table: %v", err)
	}

	// Create REGISTRATIONS table
	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		event_id INTEGER NOT NULL,
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status TEXT NOT NULL CHECK (status IN ('registered', 'canceled')),
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (event_id) REFERENCES events (id)
	)`
	_, err = DB.Exec(createRegistrationsTable)
	if err != nil {
		log.Fatalf("Could not create REGISTRATIONS table: %v", err)
	}


	
}

// func insertTestData() {
// 	// Insert sample data into USERS table
// 	insertUsers := `
// 		INSERT INTO users (username, email, password_hash, role)
// 		VALUES 
// 			('JohnDoe', 'john.doe@example.com', 'hashedpassword123', 'participant'),
// 			('JaneSmith', 'jane.smith@example.com', 'hashedpassword456', 'organizer'),
// 			('AdminUser', 'admin@example.com', 'adminpassword789', 'administrator');
// 	`
// 	_, err := DB.Exec(insertUsers)
// 	if err != nil {
// 		log.Fatalf("Could not insert into USERS table: %v", err)
// 	}

// 	// Insert sample data into EVENTS table
// 	insertEvents := `
// 		INSERT INTO events (name, description, location, start_time, end_time, status, organizer_id)
// 		VALUES 
// 			('Tech Conference', 'A conference about tech advancements.', 'New York', '2024-12-01 09:00:00', '2024-12-01 17:00:00', 'approved', 2),
// 			('Music Festival', 'An amazing music festival for all music lovers.', 'Los Angeles', '2024-12-05 14:00:00', '2024-12-05 22:00:00', 'approved', 2),
// 			('Health Symposium', 'A symposium on health and wellness.', 'San Francisco', '2024-12-10 10:00:00', '2024-12-10 16:00:00', 'pending', 2);
// 	`
// 	_, err = DB.Exec(insertEvents)
// 	if err != nil {
// 		log.Fatalf("Could not insert into EVENTS table: %v", err)
// 	}

// 	// Insert sample data into REGISTRATIONS table
// 	insertRegistrations := `
// 		INSERT INTO registrations (user_id, event_id, status)
// 		VALUES 
// 			(1, 1, 'registered'),
// 			(1, 2, 'canceled'),
// 			(3, 3, 'registered');
// 	`
// 	_, err = DB.Exec(insertRegistrations)
// 	if err != nil {
// 		log.Fatalf("Could not insert into REGISTRATIONS table: %v", err)
// 	}

// 	log.Println("Test data inserted successfully")
// }


