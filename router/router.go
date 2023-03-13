package router

import (
	"github.com/go-chi/chi"
	"z_test/handler"
)

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/home", func(home chi.Router) {
		home.Get("/", handler.HelloWorld)

		home.Route("/movie", func(mv chi.Router) {
			mv.Get("/", handler.GetMovie)
			mv.Post("/", handler.AddMovie)
			mv.Put("/", handler.UpdateMovie)
			mv.Delete("/", handler.DeleteMovie)
		})

		home.Route("/csv", func(csv chi.Router) {
			csv.Post("/", handler.ReadCSV)
		})

		home.Route("/files", func(files chi.Router) {
			files.Get("/", handler.DownloadMultipleFiles)
		})
	})

	return r
}
