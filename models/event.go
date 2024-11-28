package models

import (
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	"log"
	_ "log"
	"time"

	"abc.com/calc/db"
)


type Event struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Location    string    `json:"location"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Status      string    `json:"status"`
    Organizer   string    `json:"organizer"` // Change this to string
    CreatedAt   time.Time `json:"created_at"`
}



// func GetByName(name string) (Event, error) {
// 	var event Event
// 	query := `
// 		SELECT 
// 			e.id, 
// 			e.name, 
// 			e.description, 
// 			e.location, 
// 			e.start_time, 
// 			e.end_time, 
// 			e.status, 
// 			u.username AS organizer, 
// 			e.created_at
// 		FROM 
// 			events e
// 		INNER JOIN 
// 			users u 
// 		ON 
// 			e.organizer_id = u.id
// 		WHERE 
// 			e.name = $1 AND e.status = 'approved'
// 	`

// 	err := db.DB.QueryRow(query, name).Scan(
// 		&event.ID,
// 		&event.Name,
// 		&event.Description,
// 		&event.Location,
// 		&event.StartTime,
// 		&event.EndTime,
// 		&event.Status,
// 		&event.Organizer,
// 		&event.CreatedAt,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return Event{}, errors.New("event not found")
// 		}
// 		return Event{}, fmt.Errorf("error fetching event by name: %v", err)
// 	}

// 	return event, nil
// }

func GetByName(name string) ([]Event, error) {
	var events []Event

	query := `
		SELECT 
			e.id, 
			e.name, 
			e.description, 
			e.location, 
			e.start_time, 
			e.end_time, 
			e.status, 
			u.username AS organizer, 
			e.created_at
		FROM 
			events e
		INNER JOIN 
			users u 
		ON 
			e.organizer_id = u.id
		WHERE 
			e.name = $1 AND e.status = 'approved'
	`

	// Execute the query
	rows, err := db.DB.Query(query, name)
	if err != nil {
		return nil, fmt.Errorf("error fetching events by name: %v", err)
	}
	defer rows.Close()

	// Iterate through the rows
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.Status,
			&event.Organizer,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning event row: %v", err)
		}

		// Append the event to the slice
		events = append(events, event)
	}

	// Check for errors after iterating through rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through event rows: %v", err)
	}

	// If no events found, return a meaningful error
	if len(events) == 0 {
		return nil, errors.New("no events found with the given name")
	}

	return events, nil
}



func RegisterForEvent(userID, eventID int64) error {
	// Check if the event exists and is approved
	queryEvent := `
		SELECT COUNT(*)
		FROM events
		WHERE id = $1 AND status = 'approved'
	`
	var eventCount int
	err := db.DB.QueryRow(queryEvent, eventID).Scan(&eventCount)
	if err != nil {
		return fmt.Errorf("error checking event existence: %v", err)
	}
	if eventCount == 0 {
		return errors.New("event not found or not approved")
	}

	// Check if the user is already registered
	queryRegistration := `
		SELECT COUNT(*)
		FROM registrations
		WHERE user_id = $1 AND event_id = $2
	`
	var registrationCount int
	err = db.DB.QueryRow(queryRegistration, userID, eventID).Scan(&registrationCount)
	if err != nil {
		return fmt.Errorf("error checking registration existence: %v", err)
	}
	if registrationCount > 0 {
		return errors.New("user is already registered for this event")
	}

	// Insert the new registration
	insertQuery := `
		INSERT INTO registrations (user_id, event_id, status)
		VALUES ($1, $2, 'registered')
	`
	_, err = db.DB.Exec(insertQuery, userID, eventID)
	if err != nil {
		return fmt.Errorf("error registering for the event: %v", err)
	}

	return nil
}


func GetRegisteredEvents(userID int64) ([]Event, error) {
	var events []Event

	query := `
		SELECT 
			e.id, 
			e.name, 
			e.description, 
			e.location, 
			e.start_time, 
			e.end_time, 
			e.status, 
			u.username AS organizer, 
			e.created_at
		FROM 
			registrations r
		INNER JOIN 
			events e 
		ON 
			r.event_id = e.id
		INNER JOIN 
			users u 
		ON 
			e.organizer_id = u.id
		WHERE 
			r.user_id = $1 AND r.status = 'registered'
		ORDER BY 
			e.start_time ASC
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching registered events: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.Status,
			&event.Organizer,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning event: %v", err)
		}
		events = append(events, event)
	}

	return events, nil
}


func GetAll() ([]Event, error) {
	var events []Event
	query := `
		SELECT 
			e.id, 
			e.name, 
			e.description, 
			e.location, 
			e.start_time, 
			e.end_time, 
			e.status, 
			u.username AS organizer, 
			e.created_at
		FROM 
			events e
		INNER JOIN 
			users u 
		ON 
			e.organizer_id = u.id
		WHERE 
			e.status = 'approved'
	`
	// Log the query for debugging
	log.Printf("Executing query: %s", query)

	rows, err := db.DB.Query(query)
	if err != nil {
		// Log the error for debugging
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("error fetching approved events: %v", err)
	}
	defer rows.Close()

	// Iterate over rows
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.Status,
			&event.Organizer,
			&event.CreatedAt,
		)
		if err != nil {
			// Log scanning error for debugging
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("error scanning event row: %v", err)
		}
		events = append(events, event)
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		// Log row iteration error
		log.Printf("Error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	// Return the events slice
	return events, nil
}


func CreateEvent(event *Event, organizerID int64) error {
    query := `
        INSERT INTO events (name, description, location, start_time, end_time, status, organizer_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
        RETURNING id, created_at
    `

    // Execute the query
    err := db.DB.QueryRow(
        query,
        event.Name,
        event.Description,
        event.Location,
        event.StartTime,
        event.EndTime,
        event.Status,
        organizerID,
    ).Scan(&event.ID, &event.CreatedAt)

    if err != nil {
        log.Printf("Error creating event: %v", err) // Log the detailed error
        return fmt.Errorf("error creating event: %v", err)
    }

    return nil
}


//this is the admin function 


func GetAllAdmin() ([]Event, error) {
	var events []Event
	query := `
		SELECT 
			e.id, 
			e.name, 
			e.description, 
			e.location, 
			e.start_time, 
			e.end_time, 
			e.status, 
			u.username AS organizer, 
			e.created_at
		FROM 
			events e
		INNER JOIN 
			users u 
		ON 
			e.organizer_id = u.id
		WHERE 
			e.status = 'pending'
	`
	// Log the query for debugging
	log.Printf("Executing query: %s", query)

	rows, err := db.DB.Query(query)
	if err != nil {
		// Log the error for debugging
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("error fetching approved events: %v", err)
	}
	defer rows.Close()

	// Iterate over rows
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.Status,
			&event.Organizer,
			&event.CreatedAt,
		)
		if err != nil {
			// Log scanning error for debugging
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("error scanning event row: %v", err)
		}
		events = append(events, event)
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		// Log row iteration error
		log.Printf("Error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	// Return the events slice
	return events, nil
}


func UpdateEventStatus(eventID int64, status string) error {
	query := `UPDATE events SET status = $1 WHERE id = $2`
	result, err := db.DB.Exec(query, status, eventID)
	if err != nil {
		return fmt.Errorf("failed to update event status: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}




// func GetAll() ([]Event, error) {
// 	var events []Event
// 	query := `
// 		SELECT 
// 			e.id, 
// 			e.name, 
// 			e.description, 
// 			e.location, 
// 			e.start_time, 
// 			e.end_time, 
// 			e.status, 
// 			u.username AS organizer, 
// 			e.created_at
// 		FROM 
// 			events e
// 		INNER JOIN 
// 			users u 
// 		ON 
// 			e.organizer_id = u.id
// 		WHERE 
// 			e.status = 'approved'
// 	`

// 	rows, err := db.DB.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching approved events: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var event Event
// 		err := rows.Scan(
// 			&event.ID,
// 			&event.Name,
// 			&event.Description,
// 			&event.Location,
// 			&event.StartTime,
// 			&event.EndTime,
// 			&event.Status,
// 			&event.Organizer,
// 			&event.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning event row: %v", err)
// 		}
// 		events = append(events, event)
// 	}

// 	// Check for errors during iteration
// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error during rows iteration: %v", err)
// 	}

// 	return events, nil
// }


// func (e *Event) Save() error {
// 	query := `INSERT INTO events (name, description, location, dateTime, user_id) 
// 			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
// 	err := db.DB.QueryRow(query, e.Name, e.Description, e.Location, e.DateTime, e.UserId).Scan(&e.ID)
// 	if err != nil {
// 		return fmt.Errorf("error saving event: %v", err)
// 	}
// 	return nil
// }


// func (e *Event) Update() error {
// 	query := `UPDATE events SET name=$1, description=$2, location=$3, dateTime=$4 WHERE id=$5`
// 	_, err := db.DB.Exec(query, e.Name, e.Description, e.Location, e.DateTime, e.ID)
// 	return err
// }

// func (e *Event) DeleteIt() error {
// 	query := `DELETE FROM events WHERE id=$1`
// 	_, err := db.DB.Exec(query, e.ID)
// 	return err
// }
