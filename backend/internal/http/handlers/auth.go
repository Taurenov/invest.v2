package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/fin-helper/backend/internal/auth"
	"github.com/fin-helper/backend/internal/http/middleware"
	"github.com/fin-helper/backend/internal/repo"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	Users     repo.UserStore
	JWTSecret string
	TokenTTL  time.Duration
}

type registerReq struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "email and password (min 6) required")
		return
	}
	if req.DisplayName == "" {
		req.DisplayName = strings.Split(req.Email, "@")[0]
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal")
		return
	}

	u, err := h.Users.CreateUser(r.Context(), req.Email, hash, req.DisplayName)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "exists") {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not register")
		return
	}

	token, err := auth.IssueToken(h.JWTSecret, u.ID, u.Email, h.TokenTTL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal")
		return
	}

	writeJSON(w, map[string]any{"token": token, "user": u})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	id, hash, name, err := h.Users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal")
		return
	}
	if !auth.CheckPassword(hash, req.Password) {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := auth.IssueToken(h.JWTSecret, id, req.Email, h.TokenTTL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal")
		return
	}

	writeJSON(w, map[string]any{
		"token": token,
		"user":  map[string]any{"id": id, "email": req.Email, "display_name": name},
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := h.Users.GetByID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, map[string]any{"data": u})
}
