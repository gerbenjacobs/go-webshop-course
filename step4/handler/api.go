package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) apiProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	products, err := h.Product.ListProducts(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch products", "error", err)
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		h.logger.Error("failed to write products JSON", "error", err)
		http.Error(w, "failed to write products JSON", http.StatusInternalServerError)
	}
}

func (h *Handler) apiProductByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// validate our product ID
	productID, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		h.logger.ErrorContext(r.Context(), "couldn't convert product ID to int", "error", err)
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.Product.ShowProduct(r.Context(), productID)
	if err != nil {
		h.logger.Error("failed to fetch product", "error", err)
		http.Error(w, "failed to fetch product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		h.logger.Error("failed to write product JSON", "error", err)
		http.Error(w, "failed to write product JSON", http.StatusInternalServerError)
	}
}

func (h *Handler) apiBasket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := 1
	basket, err := h.Basket.GetBasket(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to fetch basket", "error", err)
		http.Error(w, "failed to fetch basket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(basket); err != nil {
		h.logger.Error("failed to write basket JSON", "error", err)
		http.Error(w, "failed to write basket JSON", http.StatusInternalServerError)
	}
}

func (h *Handler) apiAddToBasket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	productIDParam := r.Form.Get("product_id")

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "couldn't convert product ID to int", "error", err)
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	userID := 1
	quantity := 1
	if err := h.Basket.AddToBasket(r.Context(), userID, productID, quantity); err != nil {
		h.logger.Error("failed to add to basket", "error", err)
		http.Error(w, "failed to add to basket", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) apiRemoveFromBasket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	productIDParam := r.Form.Get("product_id")

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "couldn't convert product ID to int", "error", err)
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	userID := 1
	quantity := 1
	if err := h.Basket.RemoveFromBasket(r.Context(), userID, productID, quantity); err != nil {
		h.logger.Error("failed to remove from basket", "error", err)
		http.Error(w, "failed to remove from basket", http.StatusInternalServerError)
		return
	}
}
