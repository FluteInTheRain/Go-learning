package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	repo *ProductRepo
}

// writeJSON is a small helper to send JSON responses consistently.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// POST /products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON không hợp lệ"})
		return
	}
	if p.Name == "" || p.Price < 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Tên hoặc giá không hợp lệ"})
		return
	}
	if err := h.repo.Create(&p); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Lỗi máy chủ"})
		return
	}
	writeJSON(w, http.StatusCreated, p) // 201 with the created product (now has id)
}

// GET /products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	// read query params
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit := 100
	offset := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	products, err := h.repo.GetAll(limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Lỗi máy chủ"})
		return
	}
	writeJSON(w, http.StatusOK, products)
}

// GET /products/{id}
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id không hợp lệ"})
		return
	}
	p, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Không tìm thấy sản phẩm"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Lỗi máy chủ"})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// PUT /products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id không hợp lệ"})
		return
	}
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON không hợp lệ"})
		return
	}
	p.ID = id
	affected, err := h.repo.Update(&p)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Lỗi máy chủ"})
		return
	}
	if affected == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Không tìm thấy sản phẩm"})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// DELETE /products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id không hợp lệ"})
		return
	}
	affected, err := h.repo.Delete(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Lỗi máy chủ"})
		return
	}
	if affected == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Không tìm thấy sản phẩm"})
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204, no body
}

// COUNT NUMBER OF PRODUCTS: GET /products/count
func (h *ProductHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.repo.Count()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Loi may chu"})
		return
	}
	writeJSON(w, http.StatusOK, count)
}

func (h *ProductHandler) PartialUpdate(w http.ResponseWriter, r *http.Request) {

}
