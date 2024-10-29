package handler

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

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
		Flashes  map[string]string
		Products []app.Product
	}

	// fetch our products
	products, err := h.Product.ListProducts(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch products", "error", err)
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	flashes, err := getFlashes(r, w)
	if err != nil {
		h.logger.Warn("failed to get flashes", "error", err)
	}
	data := pageData{
		User:     false,
		Flashes:  flashes,
		Products: products,
	}
	// render the templates
	if err := tmpl.Execute(w, data); err != nil {
		h.logger.Error("failed to execute layout", "error", err)
		http.Error(w, "failed to create layout", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) productByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	productIDParam := p.ByName("id")
	h.logger.DebugContext(r.Context(), "Request received",
		"method", r.Method,
		"url", r.RequestURI,
		"product_id", productIDParam,
	)

	// validate our product ID
	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "couldn't convert product ID to int", "error", err)
		_ = storeAndSaveFlash(r, w, "warning|Invalid product ID given")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"static/layout.html",
		"static/product/product.html",
	))

	// fetch our product
	product, err := h.Product.ShowProduct(r.Context(), productID)
	switch {
	case errors.Is(err, app.ErrProductNotFound):
		h.notFound(w, r)
		return
	case err != nil:
		// an unknown error occured
		h.logger.Error("something went wrong", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// set up our page data
	type pageData struct {
		User    bool
		Flashes map[string]string
		Product app.Product
	}
	flashes, err := getFlashes(r, w)
	if err != nil {
		h.logger.Warn("failed to get flashes", "error", err)
	}
	data := pageData{
		User:    false,
		Flashes: flashes,
		Product: product,
	}

	// render the templates
	if err := tmpl.Execute(w, data); err != nil {
		h.logger.Error("failed to execute layout", "error", err)
		http.Error(w, "failed to create layout", http.StatusInternalServerError)
		return
	}
}
