package handler

import (
	"html/template"
	"net/http"

	app "github.com/gerbenjacobs/go-webshop-course"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) products(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.logger.DebugContext(r.Context(), "Request received",
		"method", r.Method,
		"url", r.RequestURI,
	)
	tmpl := template.Must(template.ParseFiles(
		"static/layout.html",
		"static/homepage.html",
	))

	type pageData struct {
		User     bool
		Products []app.Product
	}

	// fetch our products
	products, err := h.Product.ListProducts(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch products", "error", err)
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	// render the templates
	if err := tmpl.Execute(w, pageData{false, products}); err != nil {
		h.logger.Error("failed to execute layout", "error", err)
		http.Error(w, "failed to create layout", http.StatusInternalServerError)
		return
	}
}
