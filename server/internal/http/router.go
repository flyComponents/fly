package http

import (
	"net/http"

	_ "fly_server/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *Handler, photoHandler *PhotoHandler) http.Handler {

	r := chi.NewRouter()
	r.Post("/tasks", h.CreateTask)

	r.Post("/agents/{id}/heartbeat", h.Heartbeat)
	r.Post("/agents/register", h.RegisterAgent)
	r.Get("/agents", h.ListAgents)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // swagger.json URL
	))

	r.Route("/photos", func(r chi.Router) {
		r.Post("/", photoHandler.Upload)
		r.Get("/{id}", photoHandler.Get)
		r.Get("/owner/{id}", photoHandler.GetByOwnerId)
	})

	// файловый storage
	r.Get("/files/{filename}", photoHandler.ServeFile)

	return r
}
