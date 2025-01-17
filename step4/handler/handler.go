package handler

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gerbenjacobs/go-webshop-course/services"
	"github.com/julienschmidt/httprouter"
)

// Handler represents our app
// it will have dependencies
// and deal with routing
type Handler struct {
	logger *slog.Logger
	mux    http.Handler
	Dependencies
}

type Dependencies struct {
	Product services.ProductService
}

func New(logger *slog.Logger, deps Dependencies) *Handler {
	// create handler
	h := new(Handler)
	h.Dependencies = deps

	// create router
	r := httprouter.New()

	// set logger
	h.logger = logger

	// create routes
	r.GET("/", h.products)
	r.GET("/product/:id", h.productByID)

	r.NotFound = http.HandlerFunc(h.notFound)

	// set mux
	h.mux = r

	return h
}

// ServeHTTP makes it so Handler implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.logger.WarnContext(r.Context(), "Page not found",
		"method", r.Method,
		"url", r.RequestURI,
	)
	tmpl := template.Must(template.ParseFiles(
		"static/layout.html",
		"static/404.html",
	))
	w.WriteHeader(http.StatusNotFound)
	if err := tmpl.Execute(w, nil); err != nil {
		h.logger.Error("failed to execute layout", "error", err)
		http.Error(w, "failed to create layout", http.StatusInternalServerError)
		return
	}
}
