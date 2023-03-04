package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		var createTODORequest model.CreateTODORequest
		json.Unmarshal(body, &createTODORequest)
		if len(createTODORequest.Subject) > 0 {
			res, _ := h.Create(r.Context(), &createTODORequest)
			json.NewEncoder(w).Encode(res)
		} else {
			w.WriteHeader(400)
		}
	case "PUT":
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		var updateTODORequest model.UpdateTODORequest
		json.Unmarshal(body, &updateTODORequest)
		if updateTODORequest.ID != 0 && len(updateTODORequest.Subject) > 0 {
			res, err := h.Update(r.Context(), &updateTODORequest)
			if err != nil {
				w.WriteHeader(404)
			}
			json.NewEncoder(w).Encode(res)
		} else {
			w.WriteHeader(400)
		}
	default:
		w.WriteHeader(400)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, _ := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
