package internalhttp

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type handler struct {
	logger Logger
	app    Application
}

func NewHandler(logger Logger, app Application) http.Handler {
	h := &handler{
		logger: logger,
		app:    app,
	}

	r := mux.NewRouter()
	r.HandleFunc("/fill/{width}/{height}/{url:.*}", h.Fill).Methods(http.MethodGet)
	r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
	r.NotFoundHandler = http.HandlerFunc(methodNotFoundHandler)

	return r
}

func (s *handler) Fill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	width, err := strconv.ParseUint(vars["width"], 10, 32)
	if err != nil {
		badRequestResponse(w, err, s.logger)
		return
	}

	height, err := strconv.ParseUint(vars["height"], 10, 32)
	if err != nil {
		badRequestResponse(w, err, s.logger)
		return
	}

	url := vars["url"]

	bytes, err := s.app.Fill(r.Context(), r.Header, url, uint(width), uint(height))
	if err != nil {
		badRequestResponse(w, err, s.logger)
		return
	}

	successResponse(w, bytes, s.logger)
}

func methodNotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func methodNotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func badRequestResponse(w http.ResponseWriter, err error, logger Logger) {
	logger.Error(err.Error())
	http.Error(w, "502 Bad Gateway", http.StatusBadGateway)
}

func successResponse(w http.ResponseWriter, bytes []byte, logger Logger) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "image/jpeg")
	w.Header().Add("Content-Length", strconv.Itoa(len(bytes)))

	if _, err := w.Write(bytes); err != nil {
		logger.Error(err.Error())
	}
}
