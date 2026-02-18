package core

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type signupReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Message  string `json:"message"`
}

func SignupRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req signupReq
	body := io.LimitReader(r.Body, 1048576)
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if len(req.Username) <= 5 || len(req.Password) <= 5 {
		http.Error(w, `{"error":"username and password must be longer than 5 characters"}`, http.StatusBadRequest)
		return
	}

	if len(req.Message) > 500 {
		req.Message = req.Message[:500]
	}

	if Engine.UDb.CheckUserExists(req.Username) {
		http.Error(w, `{"error":"username already taken"}`, http.StatusConflict)
		return
	}

	if Engine.SRDb.HasPending(req.Username) {
		http.Error(w, `{"error":"a signup request for this username is already pending"}`, http.StatusConflict)
		return
	}

	if err := Engine.SRDb.Add(req.Username, req.Password, req.Message); err != nil {
		slog.Error("failed to create signup request", "error", err, "username", req.Username)
		http.Error(w, `{"error":"failed to submit request"}`, http.StatusInternalServerError)
		return
	}

	slog.Info("signup request submitted", "username", req.Username, "remote", r.RemoteAddr)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "pending", "message": "Your request has been submitted. An admin will review it."})
}
