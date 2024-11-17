package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"userService/domain"
	"userService/service"
	"userService/utils"

	"github.com/gorilla/mux"
)

type UsersHandler struct {
	userService service.UserService
}

func NewUsersHandler(userService service.UserService) *UsersHandler {
	return &UsersHandler{userService: userService}
}

func (h *UsersHandler) Init(r *mux.Router) {
	r.HandleFunc("/", h.GetAllUsers).Methods("GET")
	r.HandleFunc("/{id}", h.GetUserByID).Methods("GET")
	r.HandleFunc("/", h.PostUser).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	http.Handle("/", r)
}

func (h *UsersHandler) JWTAuthMiddleware(requiredRole string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(rw, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		email, role, err := utils.ValidateJWT(token)
		if err != nil || role != requiredRole {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add email to context for further handling
		r = r.WithContext(context.WithValue(r.Context(), "email", email))
		next.ServeHTTP(rw, r)
	})
}

func (h *UsersHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll()
	if err != nil || users == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	jsonResponse(users, w)
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.userService.Get(id)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	jsonResponse(user, w)
}

func (h *UsersHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err1 := h.userService.Create(&user)
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid login request", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.Email, user.UserRole)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	jsonResponse(map[string]string{"token": token}, w)
}

func jsonResponse(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
