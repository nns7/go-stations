package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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
	case "GET":
		query := r.URL.Query()
		PrevID, _ := strconv.ParseInt(query.Get("prev_id"), 10, 64)
		Size, _ := strconv.ParseInt(query.Get("size"), 10, 64)
		readTODORequest := model.ReadTODORequest{PrevID: PrevID, Size: Size}
		res, err := h.Read(r.Context(), &readTODORequest)
		if err != nil {
			w.WriteHeader(400)
		}
		json.NewEncoder(w).Encode(res)
	case "DELETE":
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		var deleteTODORequest model.DeleteTODORequest
		json.Unmarshal(body, &deleteTODORequest)
		if len(deleteTODORequest.IDs) > 0 {
			res, err := h.Delete(r.Context(), &deleteTODORequest)
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
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	res := &model.ReadTODOResponse{TODOs: make([]model.TODO, 0)}
	for _, todo := range todos {
		res.TODOs = append(res.TODOs, *todo)
	}
	return res, nil
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
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
