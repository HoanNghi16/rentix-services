package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/HoanNghi16/rentix-services/api-gateway/internal/authclient"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct{ auth *authclient.Client }

func RegisterRoutes(mux *http.ServeMux, auth *authclient.Client) {
	handler := &authHandler{auth: auth}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /auth/register", handler.register)
	mux.HandleFunc("POST /auth/login", handler.login)
}

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler *authHandler) register(w http.ResponseWriter, r *http.Request) {
	var input credentialsRequest
	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "request must contain valid email and password JSON")
		return
	}
	if strings.TrimSpace(input.Email) == "" || input.Password == "" {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "email and password are required")
		return
	}
	_, err := handler.auth.Register(r.Context(), &authclient.RegisterRequest{Email: input.Email, Password: input.Password})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "registered"})
}

func (handler *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var input credentialsRequest
	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "request must contain valid email and password JSON")
		return
	}
	if strings.TrimSpace(input.Email) == "" || input.Password == "" {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "email and password are required")
		return
	}
	response, err := handler.auth.Login(r.Context(), &authclient.LoginRequest{Email: input.Email, Password: input.Password})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request must contain a single JSON value")
	}
	return nil
}

func writeGRPCError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.InvalidArgument:
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "email and password are required")
	case codes.Unauthenticated:
		writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
	case codes.AlreadyExists:
		writeError(w, http.StatusConflict, "EMAIL_ALREADY_USED", "email is already registered")
	case codes.Unavailable, codes.DeadlineExceeded:
		writeError(w, http.StatusServiceUnavailable, "AUTH_UNAVAILABLE", "authentication service is unavailable")
	default:
		writeError(w, http.StatusBadGateway, "AUTH_REQUEST_FAILED", "authentication request failed")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"error_code": code, "message": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
