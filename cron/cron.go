package main

import (
	"chaosmanager/internal/service"
	"github.com/go-chi/chi"

	"net/http"
)

func main() {
	clientset := service.SetConfigs()
	r := chi.NewRouter()
	s := service.CronJobService{ClientSet: clientset}
	r.Route("/api", func(r chi.Router) {
		r.Post("/", s.CreateCronJob)
		r.Get("/", s.GetCronJob)
		r.Delete("/", s.DeleteCronJob)
	})
	http.ListenAndServe(":3000", r)
}
