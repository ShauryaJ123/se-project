package routes

import (
	
	"database/sql"
	"log"
	"net/http"

	"abc.com/calc/models"
	"github.com/gin-gonic/gin"
)

func viewPending(context *gin.Context) {
	// Retrieve all events
	events, err := models.GetAllAdmin()
	if err != nil {
		// Log the error if retrieval fails
		log.Printf("Error retrieving events: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events. Please try again later."})
		return
	}
	// Respond with the list of events

	if events == nil {
		events = []models.Event{} // Replace `models.Event` with the actual type of your event
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Events retrieved successfully.",
		"events":  events,
	})
}



type UpdateEventStatusRequest struct {
	EventID int64  `json:"event_id" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=approved rejected"`
}

func updateEventStatus(c *gin.Context) {
	// Parse the JSON request body
	var req UpdateEventStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body. Ensure 'event_id' and 'status' are provided."})
		return
	}

	// Update the event status using the models function
	err := models.UpdateEventStatus(req.EventID, req.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event status."})
		}
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message":  "Event status updated successfully.",
		"event_id": req.EventID,
		"status":   req.Status,
	})
}