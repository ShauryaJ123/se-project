package routes

import (
	"log"
	"net/http"
	"strconv"
	_ "strconv"

	"abc.com/calc/models"
	// "abc.com/calc/utils"
	"github.com/gin-gonic/gin"
)





func getEventByName(context *gin.Context) {
	// Get the event name from the URL parameter
	eventName := context.Param("name")

	// Check if the event name is provided
	if eventName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Event name is required. Please provide a valid name."})
		return
	}

	// Call your model to retrieve the event by name
	event, err := models.GetByName(eventName)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return the event details
	context.JSON(http.StatusOK, gin.H{
		"message": "Event retrieved successfully.",
		"event":   event,
	})
}



func getEvents(context *gin.Context) {
	// Retrieve all events
	events, err := models.GetAll()
	if err != nil {
		// Log the error if retrieval fails
		log.Printf("Error retrieving events: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events. Please try again later."})
		return
	}
	// Respond with the list of events
	context.JSON(http.StatusOK, gin.H{
		"message": "Events retrieved successfully.",
		"events":  events,
	})
}



func registerEvent(context *gin.Context) {
	// Parse JWT to get user ID
	userID, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Could not register","message":"false"})
		return
	}

	// Retrieve event ID from the path parameter
	eventIDParam := context.Param("event_id")
	if eventIDParam == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not register","message":"false"})
		return
	}

	// Convert event ID to int64
	eventID, err := strconv.ParseInt(eventIDParam, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not register","message":"false"})
		return
	}

	// Call the model to register for the event
	err = models.RegisterForEvent(userID.(int64), eventID)
	if err != nil {
		if err.Error() == "event not found or not approved" {
			context.JSON(http.StatusNotFound, gin.H{"error": "Could not register","message":"false"})
			return
		}
		if err.Error() == "user is already registered for this event" {
			context.JSON(http.StatusConflict, gin.H{"error": "You are already registered for this event.","message":"false"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register","message":"false"})
		return
	}

	// Respond with success
	context.JSON(http.StatusOK, gin.H{"message": "true"})
}







func viewRegisteredEvents(context *gin.Context) {
	// Parse JWT to get user ID
	userID, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Please log in to view your registered events."})
		return
	}

	// Call the model to fetch registered events for the user
	events, err := models.GetRegisteredEvents(userID.(int64))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch registered events. Please try again later."})
		return
	}

	// Check if there are no registered events
	if len(events) == 0 {
		context.JSON(http.StatusOK, gin.H{"message": "No registered events found."})
		return
	}

	// Respond with the list of registered events
	context.JSON(http.StatusOK, gin.H{
		"message": "Registered events retrieved successfully.",
		"events":  events,
	})
}






func createEvent(context *gin.Context) {
    var event models.Event

    // Parse JWT to get user ID
    userID, exists := context.Get("userId")
    if !exists {
        context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Please log in to create an event."})
        return
    }

    // Bind the incoming JSON to the event struct
    if err := context.ShouldBindJSON(&event); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input. Please provide valid event details."})
        return
    }

    // Validate required fields
    if event.Name == "" || event.Description == "" || event.Location == "" || event.StartTime.IsZero() || event.EndTime.IsZero() {
        context.JSON(http.StatusBadRequest, gin.H{"error": "All fields (name, description, location, start_time, end_time) are required."})
        return
    }


    // Assign the organizer using the authenticated user's ID (assumes Organizer is the user's email or username)
    organizer, err := models.GetUserById(userID.(int64))
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information. Please try again later."})
        return
    }

    event.Organizer = organizer
    event.Status = "pending" // Events start as pending and require admin approval

    // Call the model to save the event
    err = models.CreateEvent(&event, userID.(int64))
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event. Please try again later."})
        return
    }

    // Respond with success
    context.JSON(http.StatusCreated, gin.H{
        "message": "Event created successfully. It is now pending approval.",
        "event":   event,
    })
}


