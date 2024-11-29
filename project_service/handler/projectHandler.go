package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projectService/domain"
	"projectService/service"

	"github.com/gorilla/mux"
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
	r.HandleFunc("/", handler.AddProject).Methods("POST")
	r.HandleFunc("/{projectId}/addUser/{userId}", handler.AddUserToProject).Methods("PUT")
	r.HandleFunc("/user/{userId}", handler.GetProjectsByUserId).Methods("GET")
	http.Handle("/", r)
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
	var req domain.Project
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = handler.service.Create(&req)
	if err != nil {
		http.Error(w, "Unable to add project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *ProjectsHandler) AddUserToProject(w http.ResponseWriter, r *http.Request) {
	// Extract projectId and userId from request path variables
	fmt.Println("Extracting projectId and userId from request variables")
	vars := mux.Vars(r)
	projectId := vars["projectId"]
	userId := vars["userId"]
	fmt.Printf("Extracted projectId: %s, userId: %s\n", projectId, userId)

	// Call the service to add the user to the project
	fmt.Println("Calling service to add the user to the project")
	err := handler.service.AddUserToProject(projectId, userId)
	if err != nil {
		// Print the error for debugging
		fmt.Printf("Error in service.AddUserToProject: %v\n", err)

		// Return appropriate HTTP error status based on the error type
		if err.Error() == "project not found" {
			http.Error(w, "Project not found", http.StatusNotFound)
		} else if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else if err.Error() == "cannot add user: project has reached max members" {
			http.Error(w, "Project has reached maximum members", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to add user to project", http.StatusBadRequest)
		}
		return
	}

	// Respond with a success message
	fmt.Println("User successfully added to the project")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully added to the project"))
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
