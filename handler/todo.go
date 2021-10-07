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

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// methodにする
	case "POST":
		h.HandlePOST(w, r)
	case "PUT":
		h.HandleUPDATE(w, r)
	case "GET":
		h.HandleGET(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// HandlePOST is actual prosess for POST request.
func (h *TODOHandler) HandlePOST(w http.ResponseWriter, r *http.Request) {
	var todoReq model.CreateTODORequest
	err := json.NewDecoder(r.Body).Decode(&todoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if todoReq.Subject == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.Create(r.Context(), &todoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *TODOHandler) HandleUPDATE(w http.ResponseWriter, r *http.Request) {
	var todoReq model.UpdateTODORequest
	err := json.NewDecoder(r.Body).Decode(&todoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if todoReq.ID == 0 || todoReq.Subject == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := h.Update(r.Context(), &todoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *TODOHandler) HandleGET(w http.ResponseWriter, r *http.Request) {
	var (
		prevId int64
		size   int64
		err    error
	)
	if sPrevId := r.URL.Query().Get("prev_id"); sPrevId != "" {
		prevId, err = strconv.ParseInt(sPrevId, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if sSize := r.URL.Query().Get("size"); sSize != "" {
		size, err = strconv.ParseInt(sSize, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	todoReq := &model.ReadTODORequest{PrevID: prevId, Size: size}

	res, err := h.Read(r.Context(), todoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
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
