package routes

import (
	"net/http"

	"backend/handlers"
	"backend/middleware"
)

func SetupRoutes() {
	http.HandleFunc("/register", middleware.WithCORS(handlers.RegisterHandler))
	http.HandleFunc("/login", middleware.WithCORS(handlers.LoginHandler))
	http.HandleFunc("/submitComplaint", middleware.WithCORS(handlers.SubmitComplaintHandler))
	http.HandleFunc("/getAllComplaintsForUser", middleware.WithCORS(handlers.GetAllForUserHandler))
	http.HandleFunc("/getAllComplaintsForAdmin", middleware.WithCORS(handlers.GetAllForAdminHandler))
	http.HandleFunc("/viewComplaint", middleware.WithCORS(handlers.ViewComplaintHandler))
	http.HandleFunc("/resolveComplaint", middleware.WithCORS(handlers.ResolveComplaintHandler))
}
