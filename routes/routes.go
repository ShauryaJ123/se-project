package routes

import (
	"abc.com/calc/middlewares"
	"github.com/gin-gonic/gin"
)



func RegisterRoutes(server *gin.Engine) {
    // Public routes
    server.GET("/events", getEvents)               // Fetch all events
    server.POST("/signup", signup)                 // User signup
    server.POST("/login", login)    
	server.GET("/events/:name", getEventByName)               // User login
	


    // Authenticated routes
    authenticated := server.Group("/")
    authenticated.Use(middlewares.Authenticate)

    // Fetch event by name
    authenticated.POST("/events/register/:event_id", registerEvent) // Register for an event
    authenticated.GET("/events/my-registrations", viewRegisteredEvents) // View user's registered events
    // authenticated.DELETE("/events/:id/register", cancelRegistration)    // Cancel registration for an event
    authenticated.POST("/events",createEvent)

    //this function is for the admin to view pending events
    authenticated.GET("/admin/view-pending",middlewares.AuthorizeAdmin,viewPending)
    authenticated.POST("/admin/view-pending", middlewares.AuthorizeAdmin, updateEventStatus)

    //authenticated.DELETE("/events/delete/:event_id", deleteRegistration)
   
}
