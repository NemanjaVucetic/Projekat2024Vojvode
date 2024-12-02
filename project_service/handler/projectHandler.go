package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"projectService/domain"
	"projectService/service"
	"projectService/utils"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectsHandler struct {
	service *service.ProjectService
}

func NewProjectsHandler(service *service.ProjectService) *ProjectsHandler {
	return &ProjectsHandler{
		service: service,
	}
}

func (handler *ProjectsHandler) Init(r *mux.Router) {
	r.HandleFunc("/", handler.GetAllProjects).Methods("GET")
	r.HandleFunc("/{id}", handler.GetProjectByID).Methods("GET")

	//r.HandleFunc("/", handler.AddProject).Methods("POST")

	r.Handle("/", handler.JWTAuthMiddleware("Manager", http.HandlerFunc(handler.AddProject))).Methods("POST")

	r.Handle("/{projectId}/addUser/{userId}", handler.JWTAuthMiddleware("Manager", http.HandlerFunc(handler.AddUserToProject))).Methods("PUT")

	r.HandleFunc("/{projectId}/removeUser/{userId}", handler.RemoveUserFromProject).Methods("PUT")
	r.HandleFunc("/user/{userId}", handler.GetProjectsByUserId).Methods("GET")
	http.Handle("/", r)
}

func (h *ProjectsHandler) JWTAuthMiddleware(requiredRole string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(rw, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, email, role, err := utils.ValidateJWT(token) // Dodaj userID ovde
		if err != nil || role != requiredRole {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add userID and email to context for further handling
		ctx := context.WithValue(r.Context(), "userId", userID)
		ctx = context.WithValue(ctx, "email", email)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}

func (handler *ProjectsHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := handler.service.GetAll()
	if err != nil {
		http.Error(w, "Unable to fetch projects", http.StatusInternalServerError)
		return
	}

	jsonResponse(projects, w)
}

func (handler *ProjectsHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	project, err := handler.service.Get(id)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(project); err != nil {
		http.Error(w, "Unable to encode project to JSON", http.StatusInternalServerError)
	}
}

func (handler *ProjectsHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	// Extract userId from context
	userID, ok := r.Context().Value("userId").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: missing or invalid user ID", http.StatusUnauthorized)
		return
	}

	// Convert userID string to primitive.ObjectID
	managerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Decode the incoming JSON request
	var req domain.Project
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	// Set the ManagerId field to the logged-in user's ObjectID
	req.ManagerID = managerID

	// Call the service to create the project
	err = handler.service.Create(&req)
	if err != nil {
		http.Error(w, "Unable to add project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *ProjectsHandler) AddUserToProject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Extracting projectId and userId from request variables")
	vars := mux.Vars(r)
	projectId := vars["projectId"]
	userId := vars["userId"]
	fmt.Printf("Extracted projectId: %s, userId: %s\n", projectId, userId)

	loggedUserID, ok := r.Context().Value("userId").(string)
	if !ok || loggedUserID == "" {
		http.Error(w, "Unauthorized: missing or invalid user ID", http.StatusUnauthorized)
		return
	}

	fmt.Println("Calling service to add the user to the project")
	err := handler.service.AddUserToProject(projectId, userId, loggedUserID)
	if err != nil {
		fmt.Printf("Error in service.AddUserToProject: %v\n", err)

		if err.Error() == "project not found" {
			http.Error(w, "Project not found", http.StatusNotFound)
		} else if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else if err.Error() == "cannot add user: project has reached max members" {
			http.Error(w, "Project has reached maximum members", http.StatusBadRequest)
		} else if err.Error() == "user is not the project manager" {
			http.Error(w, "user is not the project manager", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to add user to project", http.StatusBadRequest)
		}
		return
	}

	fmt.Println("User successfully added to the project")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully added to the project"))
}

func (handler *ProjectsHandler) RemoveUserFromProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIdStr := vars["projectId"]
	userIdStr := vars["userId"]

	projectObjectID, err := primitive.ObjectIDFromHex(projectIdStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = handler.service.RemoveUserFromProject(projectObjectID, userObjectID)
	if err != nil {
		if err.Error() == "project not found" {
			http.Error(w, "Project not found", http.StatusNotFound)
		} else if err.Error() == "user not part of the project" {
			http.Error(w, "User not part of the project", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to remove user from project", http.StatusInternalServerError)
		}
		return
	}

	// Uspešno uklanjanje korisnika
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully removed from the project"))
}

func (handler *ProjectsHandler) GetProjectsByUserId(w http.ResponseWriter, r *http.Request) {
	// Extract userId from the request path variables
	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Call the service to fetch projects by user ID
	projects, err := handler.service.GetByUserId(userId)
	if err != nil {
		fmt.Printf("Error fetching projects for user: %v\n", err)
		http.Error(w, "Failed to fetch projects for user", http.StatusInternalServerError)
		return
	}

	// Return the projects as a JSON response
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		fmt.Printf("Error encoding projects to JSON: %v\n", err)
		http.Error(w, "Unable to encode projects to JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
