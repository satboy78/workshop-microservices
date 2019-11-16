package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type healthResponse struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

func (s Server) healthHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]healthResponse{
			"data": healthResponse{Kind: "health", Message: "everything is fine"},
		}

		b, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "error proccessing the response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}
}

type createRequest struct {
	Name      string `json:"name"`
	BelongsTo string `json:"belongs_to"`
}

func (s Server) createCounterHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "the body can't be parsed, check that is a valid json", http.StatusBadRequest)
			return
		}

		if err := s.creating.CreateCounter(ctx, req.Name, req.BelongsTo); err != nil {
			http.Error(w, "some error occurred creating process was executed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (s Server) fetchAllCountersHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		belongsTo, ok := mux.Vars(r)["belongs_to"]
		if !ok {
			http.Error(w, "the params belongs_to is required", http.StatusBadRequest)
			return
		}
		c, err := s.fetching.FetchAllCountersByUser(ctx, belongsTo)
		if err != nil {
			http.Error(w, "some error occurred fetching counters", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(c)
		if err != nil {
			http.Error(w, "error proccessing the response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}
